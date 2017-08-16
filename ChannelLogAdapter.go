package golog

// ChannelLogAdapter is a Log Adapter that logs into a buffered channel
type ChannelLogAdapter struct {
	channel chan *LogData
}

// LogData represents the data passed to the LogAdapter
type LogData struct {
	EventID      int
	ProviderID   []byte
	InstanceName string
	ProviderName string
	Message      string
	Payload      string
}

// NewChannelLogAdapter creates a new ChannelLogAdapter given the
func NewChannelLogAdapter(bufferLength int) *ChannelLogAdapter {
	return &ChannelLogAdapter{
		channel: make(chan *LogData, bufferLength),
	}
}

// GetLogChannel returns the readonly channel to receive log data
func (cal *ChannelLogAdapter) GetLogChannel() <-chan *LogData {
	return cal.channel
}

// Log will route the log data into the channel
func (cal *ChannelLogAdapter) Log(eventID int, providerID []byte, instanceName, providerName, message, payload string) (err error) {
	cal.channel <- &LogData{
		EventID:      eventID,
		ProviderID:   providerID,
		InstanceName: instanceName,
		ProviderName: providerName,
		Message:      message,
		Payload:      payload,
	}
	return
}
