package internalTypes

import "time"

type Task struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	CreatedDateTime time.Time `json:"createdDateTime"`
	Priority        int       `json:"priority"`
	IsDone          bool      `json:"isDone"`
}
