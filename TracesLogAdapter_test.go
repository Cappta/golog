package golog

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/Cappta/god"
	"github.com/Cappta/god/loki"
	testdb "github.com/erikstmartin/go-testdb"
	. "github.com/smartystreets/goconvey/convey"
)

type QueryArgs struct {
	query string
	args  []driver.Value
}

func TestTracesLogAdapter(t *testing.T) {
	Convey("Given a mocked Database", t, func() {
		queryChannel := make(chan *QueryArgs, 1)
		database, err := god.NewTestDatabase()
		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})
		testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
			queryChannel <- &QueryArgs{
				query: query,
				args:  args,
			}
			return loki.NewDriverResult(1, 1), nil
		})
		Convey("Given an Instance and Provider name", func() {
			instanceName, providerName := "TestInstance", "TestProvider"

			Convey("Given a Logger", func() {
				logger := NewLogger(NewTracesLogAdapter(database), instanceName, providerName)
				Convey("Then Logger should not be nil", func() {
					So(logger, ShouldNotBeNil)
				})
				Convey("Then Logger's InstanceName should resemble provided InstanceName", func() {
					So(logger.GetInstanceName(), ShouldResemble, instanceName)
				})
				Convey("Then Logger's ProviderName should resemble provided ProviderName", func() {
					So(logger.GetProviderName(), ShouldResemble, providerName)
				})
				Convey("Given a data format and payload", func() {
					format := "Format{data}"
					payload := map[string]interface{}{"data": "Log"}
					Convey("When logging", func() {
						err = logger.Log(1000, format, payload)
						logTime := time.Now()
						Convey("Then err should be nil", func() {
							So(err, ShouldBeNil)
						})
						Convey("Then should save in trace database", func() {
							select {
							case <-time.After(time.Second * 5):
								t.Fatal("Did not save log in the database within 5 seconds")
							case queryArgs := <-queryChannel:
								Convey("Then query should resemble configured query", func() {
									query := "INSERT INTO \"Traces\" (\"InstanceName\",\"ProviderId\",\"ProviderName\",\"EventId\",\"EventKeywords\",\"Level\",\"Opcode\",\"Task\",\"Timestamp\",\"Version\",\"FormattedMessage\",\"Payload\",\"ActivityId\",\"RelatedActivityId\",\"ProcessId\",\"ThreadId\") VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
									So(queryArgs.query, ShouldResemble, query)
								})
								Convey("Then query's InstanceName should resemble InstanceName", func() {
									So(queryArgs.args[0], ShouldResemble, instanceName)
								})
								Convey("Then query's ProviderId should resemble expected ProviderId", func() {
									So(queryArgs.args[1], ShouldResemble, []uint8{97, 8, 148, 104, 62, 243, 155, 152, 141, 51, 202, 94, 179, 218, 210, 113})
								})
								Convey("Then query's ProviderName should resemble ProviderName", func() {
									So(queryArgs.args[2], ShouldResemble, providerName)
								})
								Convey("Then query's EventID should resemble EventID", func() {
									eventID := 1000
									So(queryArgs.args[3], ShouldResemble, int64(eventID))
								})
								Convey("Then query's Timestamp should be around LogTime", func() {
									minNano := logTime.Add(-time.Second).UnixNano()
									maxNano := logTime.Add(time.Second).UnixNano()
									actualLogNano := queryArgs.args[8].(time.Time).UnixNano()
									So(actualLogNano, ShouldBeBetween, minNano, maxNano)
								})
								Convey("Then query's FormattedMessage should match expected FormattedMessage", func() {
									formattedMessage := "FormatLog"
									So(queryArgs.args[10], ShouldEqual, formattedMessage)
								})
								Convey("Then query's Payload should match expected payload", func() {
									formattedPayload := "{\"data\":\"Log\"}"
									So(queryArgs.args[11], ShouldEqual, formattedPayload)
								})
							}
						})
					})
				})
				Convey("When logging with invalid input", func() {
					err = logger.Log(1000, "{Channel}", map[string]interface{}{"Channel": make(chan int)})
					Convey("Then err should not be nil", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
