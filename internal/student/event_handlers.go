package student

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/Howard3/gosignal"
)

type eventHandlers struct {
	repo Repository
}

func NewEventHandlers(repo Repository) *eventHandlers {
	return &eventHandlers{
		repo: repo,
	}
}

// HandleNewStudentEvent is a method that handles the NewStudentEvent
// it loads the student aggregate from the repository and projects it to the database
func (eh *eventHandlers) HandleNewStudentEvent(ctx context.Context, aggregateID uint64) {
	student, err := eh.repo.loadStudent(ctx, aggregateID)
	if err != nil {
		slog.Error("failed to load student", "error", err)
		return
	}

	if err := eh.repo.upsertStudent(student); err != nil {
		slog.Error("failed to upsert student", "error", err)
		return
	}
}

// HandleUpdateStudentEvent is a method that handles the UpdateStudentEvent
// functionally the same as HandleNewStudentEvent, thus it just aliases it
func (eh *eventHandlers) HandleUpdateStudentEvent(ctx context.Context, aggID uint64) {
	eh.HandleNewStudentEvent(ctx, aggID)
}

// HandleGenerateCodeEvent is a method that handles the GenerateCodeEvent
func (eh *eventHandlers) HandleGenerateCodeEvent(ctx context.Context, aggID uint64) {
	student, err := eh.repo.loadStudent(ctx, aggID)
	if err != nil {
		slog.Error("failed to load student", "error", err)
		return
	}

	code := student.data.CodeUniqueId
	if len(code) == 0 {
		slog.Error("code is empty")
		return
	}

	if err := eh.repo.insertStudentCode(ctx, aggID, code); err != nil {
		slog.Error("failed to insert student code", "error", err)
		return
	}
}

// routeEvent is a method that routes an event to the appropriate handler
func (eh *eventHandlers) routeEvent(ctx context.Context, evt *gosignal.Event) {
	id, err := strconv.ParseUint(evt.AggregateID, 10, 64)
	if err != nil {
		slog.Error("failed to parse aggregate ID", "error", err)
		return
	}

	switch evt.Type {
	case EVENT_ADD_STUDENT:
		eh.HandleNewStudentEvent(ctx, id)
	case EVENT_UPDATE_STUDENT, EVENT_ENROLL_STUDENT, EVENT_UNENROLL_STUDENT, EVENT_SET_STUDENT_STATUS:
		eh.HandleUpdateStudentEvent(ctx, id)
	case EVENT_SET_LOOKUP_CODE:
		eh.HandleGenerateCodeEvent(ctx, id)
	}
}
