package domain

import "fmt"

type TaskStatus int

const (
	TaskStatusCreated   TaskStatus = 1
	TaskStatusCompleted TaskStatus = 2
	TaskStatusFailed    TaskStatus = 3
)

func (ts TaskStatus) String() string {
	switch ts {
	case TaskStatusCreated:
		return "created"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", ts)
	}
}

func (ts TaskStatus) IsValid() bool {
	return ts == TaskStatusCreated || ts == TaskStatusCompleted || ts == TaskStatusFailed
}

func TaskStatusFromInt(value int) (TaskStatus, error) {
	status := TaskStatus(value)
	if !status.IsValid() {
		return 0, fmt.Errorf("invalid task status: %d", value)
	}
	return status, nil
}

func TaskStatusFromString(value string) (TaskStatus, error) {
	switch value {
	case "created":
		return TaskStatusCreated, nil
	case "completed":
		return TaskStatusCompleted, nil
	case "failed":
		return TaskStatusFailed, nil
	default:
		return 0, fmt.Errorf("invalid task status string: %s", value)
	}
}
