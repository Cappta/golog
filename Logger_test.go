package golog

import (
	"errors"
	"testing"

	"github.com/Cappta/debugo"
	. "github.com/Cappta/gohelpconvey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	Convey("Given a ChannelLogAdapter", t, func() {
		channelLogAdapter := NewChannelLogAdapter(1)
		Convey("Then ChannelLogAdapter not be nil", func() {
			So(channelLogAdapter, ShouldNotBeNil)
		})

		Convey("Given an Instance and Provider name", func() {
			instanceName, providerName := "TestInstance", "TestProvider"

			Convey("Given a Logger", func() {
				logger := NewLogger(channelLogAdapter, instanceName, providerName)

				conveyLog := func(logFunc func() error, eventID int, message, payload string) {
					err := logFunc()
					Convey("Then err should be nil", func() {
						So(err, ShouldBeNil)
					})

					logData := <-channelLogAdapter.GetLogChannel()
					Convey("Then InstanceName should resemble InstanceName", func() {
						So(logData.InstanceName, ShouldResemble, instanceName)
					})
					Convey("Then query's ProviderId should resemble expected ProviderId", func() {
						providerID := []uint8{97, 8, 148, 104, 62, 243, 155, 152, 141, 51, 202, 94, 179, 218, 210, 113}
						So(logData.ProviderID, ShouldResemble, providerID)
					})
					Convey("Then query's ProviderName should resemble ProviderName", func() {
						So(logData.ProviderName, ShouldResemble, providerName)
					})
					Convey("Then query's EventID should resemble EventID", func() {
						So(logData.EventID, ShouldResemble, eventID)
					})
					Convey("Then query's FormattedMessage should match expected FormattedMessage", func() {
						So(logData.Message, ShouldMatch, message)
					})
					Convey("Then query's Payload should match expected payload", func() {
						So(logData.Payload, ShouldMatch, payload)
					})
				}

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
						conveyLog(
							func() error { return logger.Log(1000, format, payload) },
							1000,
							"FormatLog",
							"{\"data\":\"Log\"}",
						)
					})
					Convey("Given a message", func() {
						message := "Cappta melhor MAE :D"
						Convey("When logging info", func() {
							conveyLog(func() error { return logger.Info(message) },
								1000,
								"Host: CAPPDESK-0103; Message: Cappta melhor MAE :D",
								"{\"hostName\":\"CAPPDESK-0103\",\"message\":\"Cappta melhor MAE :D\"}",
							)
						})
					})
					Convey("Given an error", func() {
						err := errors.New("Filho da MAE")
						Convey("When logging warning", func() {
							conveyLog(func() error { return logger.Warning(err) },
								2000,
								"Host: .*; Operation: github.com/Cappta/golog.TestLogger.func.*; FileName: .*/github.com/Cappta/golog/Logger_test.go: LineNumber: \\d+; Exception: Filho da MAE",
								"{\"err\":\"Filho da MAE\",\"fileName\":\".*/golog/Logger_test.go\",\"host\":\".*\",\"lineNumber\":\\d+,\"operation\":\"github.com/Cappta/golog.TestLogger.func.*\"}",
							)
						})
						Convey("When logging error", func() {
							conveyLog(func() error { return logger.Error(err) },
								3000,
								"Host: .*; Operation: github.com/Cappta/golog.TestLogger.func.*; FileName: .*/github.com/Cappta/golog/Logger_test.go: LineNumber: \\d+; Exception: Filho da MAE",
								"{\"err\":\"Filho da MAE\",\"fileName\":\".*/golog/Logger_test.go\",\"host\":\".*\",\"lineNumber\":\\d+,\"operation\":\"github.com/Cappta/golog.TestLogger.func.*\"}",
							)
						})
						Convey("When mocking os.Hostname to return an error", func() {
							returnedError := errors.New("FODEU O HOSTNAME")
							osHostname = func() (name string, err error) {
								return "", returnedError
							}
							Convey("When logging warning", func() {
								err = logger.Warning(err)
								Convey("Then error should resemble returned error", func() {
									So(err, ShouldResemble, returnedError)
								})
							})
							Convey("When logging error", func() {
								err = logger.Error(err)
								Convey("Then error should resemble returned error", func() {
									So(err, ShouldResemble, returnedError)
								})
							})
						})
						Convey("When mocking Debug.GetCaller to return an error", func() {
							returnedError := errors.New("FODEU O DEBUG")
							debugGetCaller = func(depth int) (stackCall *debugo.StackCall, err error) {
								return nil, returnedError
							}
							Convey("When logging warning", func() {
								err = logger.Warning(err)
								Convey("Then error should resemble returned error", func() {
									So(err, ShouldResemble, returnedError)
								})
							})
							Convey("When logging error", func() {
								err = logger.Error(err)
								Convey("Then error should resemble returned error", func() {
									So(err, ShouldResemble, returnedError)
								})
							})
						})
					})
				})
				Convey("When logging with invalid input", func() {
					err := logger.Log(1000, "{Channel}", map[string]interface{}{"Channel": make(chan int)})
					Convey("Then err should not be nil", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
