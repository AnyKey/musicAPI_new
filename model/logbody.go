package model

import "time"

type LogBody struct {
	Name   string    `json:"name"`
	Action string    `json:"action"`
	Time   time.Time `json:"time"`
}
