package models

import "fmt"

type X10Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *X10Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("X10Error [%d]: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("X10Error [%d]: %s", e.Code, e.Message)
}
