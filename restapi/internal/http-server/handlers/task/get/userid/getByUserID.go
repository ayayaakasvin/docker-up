package get

import (
	"errors"
	"log/slog"
	"net/http"

	"restapi/internal/errorset"
	"restapi/internal/http-server/handlers/task/get"
	helper "restapi/internal/lib/helperfunctions"
	"restapi/internal/lib/sl"
	"restapi/internal/models"
	"restapi/internal/models/data"
	"restapi/internal/models/response"

	"github.com/gin-gonic/gin"
)

func GetUserTasksHandler(log *slog.Logger, tg get.TasksGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load logger with necessary data
		const op = "handlers.task.get.GetUserTasksHandler"
		helper.LoadLogger(&log, c, op)

		// fetch ID param
		userId := helper.GetIDFromParams(c, helper.UserIDKey)
		if userId == -1 {
			response.Error(c, http.StatusBadRequest, errorset.ErrBindRequest)
			return
		}

		log.Info("decoded request", slog.Any(helper.UserIDKey, userId))

		// action with db
		tasksSlice, err := tg.GetTasksByUserID(userId)
		if err != nil || len(tasksSlice) == 0 {
			handleGettingTasksError(c, log, err, tasksSlice, userId)
			return
		}

		var data data.Data = data.NewData()
		data[helper.TasksKey] = tasksSlice

		log.Info("tasks succesfully passed", slog.Int64(helper.UserIDKey, userId))
		response.Ok(c, http.StatusOK, data)
	}
}

func handleGettingTasksError(c *gin.Context, log *slog.Logger, err error, tasksSlice []*models.Task, userId int64)  {
	if errors.Is(err, errorset.ErrUserNotFound) {
		log.Error(err.Error(), sl.Err(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	} 
	
	if len(tasksSlice) == 0 {
		log.Warn("no tasks found for user", slog.Int64(helper.UserIDKey, userId))
		response.Error(c, http.StatusNotFound, errorset.ErrTaskNotFound.Error())
		return
	}
}