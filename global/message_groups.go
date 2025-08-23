package global

type MessageGroup string

const (
	//Internal
	UnknownMessageGroup       = "unknown"
	InternalIssueMessageGroup = "Internal issue"

	//Permission
	PermissionDeniedMessageGroup = "Permission denied"

	//Validation
	InvalidInputMessageGroup = "Invalid input"
)
