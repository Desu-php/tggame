package controllers

import (
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type TaskController struct {
	logger         *logrus.Logger
	taskRepository repository.TaskRepository
	taskService    *services.TaskService
}

func NewTaskController(
	logger *logrus.Logger,
	taskRepository repository.TaskRepository,
	taskService *services.TaskService,
) *TaskController {
	return &TaskController{
		logger:         logger,
		taskRepository: taskRepository,
		taskService:    taskService,
	}
}

func (c *TaskController) GetAll(ctx *gin.Context) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return
	}

	tasks, err := c.taskRepository.GetAll(ctx, user)

	if err != nil {
		c.logger.WithError(err).Error("TaskController::GetAll")
		responses.ServerErrorResponse(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

func (c *TaskController) ClickLink(ctx *gin.Context) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return
	}

	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	taskID := uint(idUint64)

	ok, err = c.taskService.ProgressOnTimeTask(ctx, &services.ProgressOneTimeDto{
		User:     user,
		TaskID:   taskID,
		TaskType: models.TaskTypeClickLink,
	})

	if err != nil {
		c.logger.WithError(err).Error("TaskController::ClickLink")
		responses.ServerErrorResponse(ctx)
		return
	}

	if ok == false {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  errs.TaskAlreadyCompletedCode,
			"error": "Task already completed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *TaskController) ReceiveReward(ctx *gin.Context) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return
	}

	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id := uint(idUint64)

	userTask, err := c.taskService.ReceiveReward(ctx, id, user)

	if err != nil {
		c.logger.WithError(err).Error("TaskController::ReceiveReward")
		responses.ServerErrorResponse(ctx)
		return
	}

	if userTask == nil {
		responses.NotFound(ctx)
		return
	}

	if userTask.CompletedAt == nil {
		ctx.JSON(http.StatusOK, gin.H{"amount": 0, "completed": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"amount": userTask.Task.Amount, "completed": true})
}
