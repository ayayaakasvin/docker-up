package task

import (
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

// SaveTask implements TaskHandlers.
func (t TaskHandler) SaveTask(c *gin.Context) {
	// load logger with necessary data
	const op = "handlers.task.TaskHandler.SaveTask"
	logger := helper.LoadLogger(t.log, c, op)

	// fetch ID param
	userID := helper.GetIDFromParams(c, helper.UserIDKey)
	if userID == -1 {
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	// bind request
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(errorset.ErrBindRequest, sl.Err(err))
		response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
		return
	}

	logger.Info("decoded request", slog.Any(helper.ReqKey, req))

	// action with db
	taskId, err := t.db.SaveTask(userID, req.TaskContent)
	if err != nil {
		handleSavingTaskError(c, logger, err, taskId)
		return
	}

	var data data.Data = data.NewData()
	data[helper.UserIDKey] = userID
	data[helper.TaskIDKey] = taskId

	logger.Info("task saved successfully", slog.Int64(helper.UserIDKey, userID))
	logger.Info("task saved successfully", slog.Int64(helper.TaskIDKey, taskId))

	response.Ok(c, http.StatusCreated, data)
}

func handleSavingTaskError(c *gin.Context, log *slog.Logger, err error, taskId int64) {
	log.Error("failed to save task", sl.Err(err))
	if err == errorset.ErrUserNotFound {
		response.Error(c, http.StatusConflict, errorset.ErrUserNotFound.Error())
		return
	} else if taskId == 0 {
		log.Error("unexpected task ID = 0 after saving task")
		response.Error(c, http.StatusInternalServerError, "unexpected server error")
		return
	}
}
