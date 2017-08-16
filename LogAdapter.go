package golog

// LogAdapter represents a logger adapter
type LogAdapter interface {
	Log(eventID int, providerID []byte, instanceName, providerName, message, payload string) (err error)
}
