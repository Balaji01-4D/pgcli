package database

import (
	"fmt"
	"time"
)

// exec → returns status/affected rows (INSERT, UPDATE, CREATE, etc.)

type ExecResult struct {
	RowsAffected int64
	Status       string
	Duration     time.Duration
}

func (e *ExecResult) GetType() string {
	return "EXEC"
}

func (e *ExecResult) Render() string {
	return fmt.Sprintf(
		"Status: %s\nRows Affected: %d\nDuration: %s",
		e.Status,
		e.RowsAffected,
		e.Duration.String(),
	)
}