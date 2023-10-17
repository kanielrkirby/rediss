package rerror

import (
	"fmt"
)

// RedisError: Code, Message, and Args (args are string formatting args).
type RedisError struct {
	Code    string
	Message string
	Args    []interface{}
}

// RedisErrorCode: Code and Message.
type RedisErrorCode struct {
	Code    string
	Message string
}

// DefineError: Define a new error code and message.
func DefineError(code string, message string) *RedisErrorCode {
	return &RedisErrorCode{Code: code, Message: message}
}

// A map of RedisErrorCodes to their messages.
var (
	ERR_OK                  = DefineError("OK", "OK")
	ERR_CANCELLED           = DefineError("CANCELLED", "Operation cancelled")
	ERR_UNKNOWN             = DefineError("UNKNOWN", "Unknown error")
	ERR_INVALID_ARGUMENT    = DefineError("INVALID_ARGUMENT", "Invalid argument")
	ERR_DEADLINE_EXCEEDED   = DefineError("DEADLINE_EXCEEDED", "Deadline exceeded")
	ERR_NOT_FOUND           = DefineError("NOT_FOUND", "Key not found")
	ERR_ALREADY_EXISTS      = DefineError("ALREADY_EXISTS", "Key already exists")
	ERR_PERMISSION_DENIED   = DefineError("PERMISSION_DENIED", "Permission denied")
	ERR_RESOURCE_EXHAUSTED  = DefineError("RESOURCE_EXHAUSTED", "Resource exhausted")
	ERR_FAILED_PRECONDITION = DefineError("FAILED_PRECONDITION", "Failed precondition")
	ERR_ABORTED             = DefineError("ABORTED", "Aborted")
	ERR_OUT_OF_RANGE        = DefineError("OUT_OF_RANGE", "Out of range")
	ERR_UNIMPLEMENTED       = DefineError("UNIMPLEMENTED", "Command not implemented yet, but is planned.")
	ERR_INTERNAL            = DefineError("INTERNAL", "Internal error")
	ERR_UNAVAILABLE         = DefineError("UNAVAILABLE", "Service unavailable")
	ERR_DATA_LOSS           = DefineError("DATA_LOSS", "Data loss")
	ERR_UNAUTHENTICATED     = DefineError("UNAUTHENTICATED", "Unauthenticated")
	ERR_UNKNOWN_SUBCOMMAND  = DefineError("UNKNOWN_SUBCOMMAND", "Unknown subcommand '%s'. Try %s HELP.")
)

// Temporary boolean for debugging
var isDebug bool = true

// Error: Return the error message, formatted with Sprintf.
func (e *RedisError) Error() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

// New: Provide args to format the error message.
func New(Code *RedisErrorCode, args ...interface{}) error {
	return &RedisError{Code: Code.Code, Message: Code.Message, Args: args}
}
