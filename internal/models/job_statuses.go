package models

const (
	// CreatedStatus is a status when job is created
	CreatedStatus int = iota + 1

	// ProcessingStatus is a status when job is created
	ProcessingStatus

	// FailedStatus is a status when job is failed
	FailedStatus

	// NoAttemptsLeftStatus is a status when job doesn't have any attempts left
	NoAttemptsLeftStatus

	// DoneStatus is a status when job is done
	DoneStatus
)
