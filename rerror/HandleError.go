package rerror

import (
  "fmt"
)

type RedisError struct {
  Number int
  GrpcDescription string
  Message string
}

var KnownStatus = map[string]RedisError{
  "OK": {Number: 0, GrpcDescription: "Not an error; returned on success.", Message: "OK"},
  "CANCELLED": {Number: 1, GrpcDescription: "The operation was cancelled, typically by the caller.", Message: "Operation cancelled"},
  "UNKNOWN": {Number: 2, GrpcDescription: "Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.", Message: "Unknown error"},
  "INVALID_ARGUMENT": {Number: 3, GrpcDescription: "The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).", Message: "Invalid argument"},
  "DEADLINE_EXCEEDED": {Number: 4, GrpcDescription: "The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long", Message: "Deadline exceeded"},
  "NOT_FOUND": {Number: 5, GrpcDescription: "Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used. This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.", Message: "Key not found"},
  "ALREADY_EXISTS": {Number: 6, GrpcDescription: "The entity that a client attempted to create (e.g., file or directory) already exists.", Message: "Key already exists"},
  "PERMISSION_DENIED": {Number: 7, GrpcDescription: "The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some resource (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.", Message: "Permission denied"},
  "RESOURCE_EXHAUSTED": {Number: 8, GrpcDescription: "Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.", Message: "Resource exhausted"},
  "FAILED_PRECONDITION": {Number: 9, GrpcDescription: "The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an \"rmdir\" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.", Message: "Failed precondition"},
  "ABORTED": {Number: 10, GrpcDescription: "The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.", Message: "Aborted"},
  "OUT_OF_RANGE": {Number: 11, GrpcDescription: "The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.", Message: "Out of range"},
  "UNIMPLEMENTED": {Number: 12, GrpcDescription: "The operation is not implemented or is not supported/enabled in this service.", Message: "Command not implemented yet, but is planned."},
  "INTERNAL": {Number: 13, GrpcDescription: "Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.", Message: "Internal error"},
  "UNAVAILABLE": {Number: 14, GrpcDescription: "The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.", Message: "Service unavailable"},
  "DATA_LOSS": {Number: 15, GrpcDescription: "Unrecoverable data loss or corruption.", Message: "Data loss"},
  "UNAUTHENTICATED": {Number: 16, GrpcDescription: "The request does not have valid authentication credentials for the operation.", Message: "Unauthenticated"},
  "UNKNOWN_SUBCOMMAND": {Number: 17, GrpcDescription: "The subcommand is not known.", Message: "Unknown subcommand '%s'. Try %s HELP."},
}

func NewRedisError(code string, message string) RedisError {
  return RedisError{
    Number: KnownStatus[code].Number,
    GrpcDescription: KnownStatus[code].GrpcDescription,
    Message: message,
  }
}

func ArgumentCount(command string) string {
  return fmt.Sprintf("ERR wrong number of arguments for '%s' command", command)
}
