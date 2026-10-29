package enum

type UpdateStatus string

const (
	UpdateStatusPending UpdateStatus = "pending"
	UpdateStatusDone    UpdateStatus = "done"
	UpdateStatusFailed  UpdateStatus = "failed"
)
