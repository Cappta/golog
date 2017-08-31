package golog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Cappta/gohelpmath"
)

// Logger represents a Logger
type Logger struct {
	adapter      LogAdapter
	instanceName string
	providerName string
	providerID   []byte
}

const idLength = 16

// NewLogger creates a new Logger with the specified adapter, instance and provider names
func NewLogger(adapter LogAdapter, instanceName, providerName string) *Logger {
	return &Logger{
		adapter:      adapter,
		instanceName: instanceName,
		providerName: providerName,
		providerID:   gohelpmath.Hash([]byte(providerName), idLength),
	}
}

// GetInstanceName returns the instance name of this logger
func (logger *Logger) GetInstanceName() string {
	return logger.instanceName
}

// GetProviderName returns the provider name of this logger
func (logger *Logger) GetProviderName() string {
	return logger.providerName
}

// Log will log the provided event with the specified format and payload
func (logger *Logger) Log(eventID int, format string, payload map[string]interface{}) (err error) {
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}

	message := format
	for key, value := range payload {
		keyWrapper := fmt.Sprintf("{%s}", key)
		stringValue := fmt.Sprintf("%v", value)
		message = strings.Replace(message, keyWrapper, stringValue, -1)
	}
	return logger.adapter.Log(eventID, logger.providerID, logger.instanceName, logger.providerName, message, string(marshalledPayload))
}

// Info logs an info message
func (logger *Logger) Info(message string) error {
	hostName, _ := os.Hostname()
	return logger.Log(1000,
		"Host: {hostName}; Message: {message}",
		map[string]interface{}{"hostName": hostName, "message": message},
	)
}

// Warning logs an error as a warning
func (logger *Logger) Warning(err error) (returnError error) {
	stackCall, returnError := debugGetCaller(0)
	if returnError != nil {
		return
	}
	operation := stackCall.Func.Name()
	hostName, returnError := osHostname()
	if returnError != nil {
		return
	}
	return logger.Log(2000,
		"Host: {host}; Operation: {operation}; FileName: {fileName}: LineNumber: {lineNumber}; Exception: {err}",
		map[string]interface{}{"host": hostName, "operation": operation, "fileName": stackCall.File, "lineNumber": stackCall.Line, "err": err.Error()},
	)
}

// Error logs an error as a failure
func (logger *Logger) Error(err error) (returnError error) {
	stackCall, returnError := debugGetCaller(0)
	if returnError != nil {
		return
	}
	operation := stackCall.Func.Name()
	hostName, returnError := osHostname()
	if returnError != nil {
		return
	}
	return logger.Log(3000,
		"Host: {host}; Operation: {operation}; FileName: {fileName}: LineNumber: {lineNumber}; Exception: {err}",
		map[string]interface{}{"host": hostName, "operation": operation, "fileName": stackCall.File, "lineNumber": stackCall.Line, "err": err.Error()},
	)
}
