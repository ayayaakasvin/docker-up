package save

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

type Request struct {
	TaskContent string `json:"taskContent" binding:"required"`
}

type TaskSaver interface {
	SaveTask(userId int64, content string) (int64, error)
}

func SaveTaskHandler(log *slog.Logger, ts TaskSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.task.save.SaveTaskHandler"
		helper.LoadLogger(&log, c, op)

		// fetch ID param
		userID := helper.GetIDFromParams(c, helper.UserIDKey)
		if userID == -1 {
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		// bind request
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(errorset.ErrBindRequest, sl.Err(err))
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Any(helper.ReqKey, req))

		// action with db
		taskId, err := ts.SaveTask(userID, req.TaskContent)
		if err != nil {
			handleSavingTaskError(c, log, err, taskId)
			return
		}

		var data data.Data = data.NewData()
		data[helper.UserIDKey] = userID
		data[helper.TaskIDKey] = taskId

		log.Info("task saved successfully", slog.Int64(helper.UserIDKey, userID))
		log.Info("task saved successfully", slog.Int64(helper.TaskIDKey, taskId))

		response.Ok(c,http.StatusCreated, data)
	}
}

func handleSavingTaskError (c *gin.Context, log *slog.Logger, err error, taskId int64)  {
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