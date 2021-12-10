package models

type Status int64

const (
	Pending Status = 0
	Rejected Status = 1
	Accepted Status = 2
)

func (s Status) String() string {
	switch s {
	case 0:
		return "pending"
	case 1:
		return "rejected"
	case 2:
		return "accepted"
	default:
		return "unknown"
	}
}