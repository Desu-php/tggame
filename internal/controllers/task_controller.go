package controllers

import (
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	"example.com/v2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TaskController struct {
	logger         *logrus.Logger
	taskRepository repository.TaskRepository
}

func NewTaskController(
	logger *logrus.Logger,
	taskRepository repository.TaskRepository,
) *TaskController {
	return &TaskController{
		logger:         logger,
		taskRepository: taskRepository,
	}
}

func (c *TaskController) GetAll(ctx *gin.Context) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return
	}

	tasks, err := c.taskRepository.GetAll(ctx, user)

	if err != nil {
		c.logger.WithError(err).Error("failed to fetch tasks")
		responses.ServerErrorResponse(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}
