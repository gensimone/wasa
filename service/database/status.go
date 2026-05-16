package database

type Status string

const (
	Sent     Status = "sent"
	Received Status = "received"
	Read     Status = "read"
)

var validStatus = map[Status]struct{}{
	Sent:     {},
	Received: {},
	Read:     {},
}

func IsValidStatus(status Status) bool {
	_, ok := validStatus[status]
	return ok
}
