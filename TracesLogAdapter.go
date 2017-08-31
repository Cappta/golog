package golog

import (
	"time"

	"github.com/Cappta/god"
	"github.com/Cappta/god/CapptaLog"
)

// TracesLogAdapter represents a log adapter which will log in the Traces table
type TracesLogAdapter struct {
	database *god.Database
}

// NewTracesLogAdapter creates a new TracesLogAdapter
func NewTracesLogAdapter(database *god.Database) *TracesLogAdapter {
	return &TracesLogAdapter{
		database: database,
	}
}

// Log will log the specified event into the Traces table
func (logger *TracesLogAdapter) Log(eventID int, providerID []byte, instanceName, providerName, message, payload string) (err error) {
	traces := &CapptaLog.Traces{
		InstanceName:     instanceName,
		ProviderID:       providerID,
		ProviderName:     providerName,
		EventID:          eventID,
		Timestamp:        time.Now(),
		FormattedMessage: &message,
		Payload:          &payload,
	}

	db := logger.database.Save(traces)
	return db.Error
}
