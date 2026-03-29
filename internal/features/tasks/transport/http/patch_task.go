package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       types.Nullable[string] `json:"title"`
	Description types.Nullable[string] `json:"description"`
	Completed   types.Nullable[bool]   `json:"completed"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("required field `title` cant be null")
		}
		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("required field `title` length must be between 1 and 100")
		}
	}
	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("field `description` length must be between 1 and 1000")
			}
		}
	}
	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("required field `completed` cant be null")
		}
	}

	return nil
}

// PatchTask godoc
// @Summary Частичное обновление задачи
// @Description Частично обновляет данные задачи по идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body PatchTaskRequest true "Тело запроса для частичного обновления задачи"
// @Success 200 {object} PatchTaskResponse "Задача успешно обновлена"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id")
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate  request")
	}

	taskPatch := taskPatchFromRequest(request)
	taskDomain, err := h.tasksService.PatchTask(ctx, taskId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}
	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
