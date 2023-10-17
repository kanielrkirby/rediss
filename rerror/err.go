package rerror

import (
	"fmt"
  "runtime"
  "strings"
)

// RedisError represents a custom error code and message, with args to format the message using fmt.Sprintf.
type RedisError struct {
	Code    string
	Message string
}

// defineError returns a standard RedisError with the given code and message.
func defineError(code string, message string) *RedisError {
  return &RedisError{Code: code, Message: message}
}

var (
  // ErrOk is used to indicate that the operation was successful.
	ErrOk                  = defineError("OK", "OK")
  // ErrCancelled is used to indicate that the operation was cancelled.
	ErrCancelled           = defineError("CANCELLED", "Operation cancelled")
  // ErrUnknown is used to indicate that the operation failed for an unknown reason.
	ErrUnknown             = defineError("UNKNOWN", "Unknown error")
  // ErrInvalidArgument is used to indicate that the operation was given an invalid argument.
	ErrInvalidArgument    = defineError("INVALID_ARGUMENT", "Invalid argument")
  // ErrDeadlineExceeded is used to indicate that the operation exceeded its deadline.
	ErrDeadlineExceeded   = defineError("DEADLINE_EXCEEDED", "Deadline exceeded")
  // ErrNotFound is used to indicate that the requested entity was not found.
	ErrNotFound           = defineError("NOT_FOUND", "Key not found")
  // ErrAlreadyExists is used to indicate that the entity that a caller attempted to create already exists.
	ErrAlreadyExists      = defineError("ALREADY_EXISTS", "Key already exists")
  // ErrPermissionDenied is used to indicate that the caller does not have permission to execute the specified operation.
	ErrPermissionDenied   = defineError("PERMISSION_DENIED", "Permission denied")
  // ErrResourceExhausted is used to indicate that some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	ErrResourceExhausted  = defineError("RESOURCE_EXHAUSTED", "Resource exhausted")
  // ErrFailedPrecondition is used to indicate that the operation was rejected because the system is not in a state required for the operation's execution.
	ErrFailedPrecondition = defineError("FAILED_PRECONDITION", "Failed precondition")
  // ErrAborted is used to indicate that the operation was aborted, typically due to a concurrency issue like sequencer check failures, transaction aborts, etc.
	ErrAborted             = defineError("ABORTED", "Aborted")
  // ErrOutOfRange is used to indicate that the operation was attempted past the valid range.
	ErrOutOfRange        = defineError("OUT_OF_RANGE", "Out of range")
  // ErrUnimplemented is used to indicate that the operation is not implemented or is not supported/enabled in this service.
	ErrUnimplemented       = defineError("UNIMPLEMENTED", "Command not implemented yet, but is planned.")
  // ErrInternal is used to indicate that an internal error occurred.
	ErrInternal            = defineError("INTERNAL", "Internal error")
  // ErrUnavailable is used to indicate that the service is currently unavailable.
	ErrUnavailable         = defineError("UNAVAILABLE", "Service unavailable")
  // ErrDataLoss is used to indicate that unrecoverable data loss or corruption occurred.
	ErrDataLoss           = defineError("DATA_LOSS", "Data loss")
  // ErrUnauthenticated is used to indicate that the request does not have valid authentication credentials for the operation.
	ErrUnauthenticated     = defineError("UNAUTHENTICATED", "Unauthenticated")
  // ErrUnknownSubcommand is used to indicate that the subcommand is not known. Provide "Subcommand" and "Command" to format the message.
	ErrUnknownSubcommand  = defineError("UNKNOWN_SUBCOMMAND", "Unknown subcommand '%s'. Try %s HELP.")
  // ErrWrongNumberOfArguments is used to indicate that the wrong number of arguments were provided. Provide "Subcommand" and "Command" to format the message.
  ErrWrongNumberOfArguments = defineError("WRONG_NUMBER_OF_ARGUMENTS", "Wrong number of arguments for '%s'. Try %s HELP.")
  // ErrUnknownType is used to indicate that the type of the value is unknown. Provide "Type" to format the message.
  ErrUnknownType = defineError("UNKNOWN_TYPE", "Unknown type '%v'.")
  // ErrWrap is used when no other error code is appropriate. Provide the error and a custom "Message" to format the message.
  ErrWrap = func(err error) *RedisError {
    return defineError("WRAP", "%s:" + err.Error())
  }
)

// DEBUG is used to determine whether to show the function and line number of the error.
var DEBUG bool = true

// Error formats the error message using fmt.Sprintf and returns it.
func (e *RedisError) Error() string {
  if !DEBUG {
    return e.Message
  }
 builder := strings.Builder{}
 for i := 1; i < 10; i++ {
   pc, file, line, _ := runtime.Caller(i)
   // if it exists
   if pc == 0 {
     break
   }
   fn := runtime.FuncForPC(pc).Name()
   builder.WriteString(fmt.Sprintf("%s:%d in %s\n", file, line, fn))
 }
  builder.WriteString(e.Message)
  return builder.String()
}

// Format formats the error message using fmt.Sprintf and returns it.
func (e *RedisError) Format(args ...interface{}) *RedisError {
  if (args == nil || len(args) == 0) {
    return e
  }
  e.Message = fmt.Sprintf(e.Message, args...)
  return e
}

// FormatAndError is shorthand syntax for .Format(...args).Error()
func (e *RedisError) FormatAndError(args ...interface{}) string {
  return e.Format(args...).Error()
}

// Usage: to create a new error, use rerror.New(rerror.ErrUnimplemented, "some", "args")
