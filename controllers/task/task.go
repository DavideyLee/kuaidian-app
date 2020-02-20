package taskcontrollers

import (
	"kuaidian-app/controllers"

	"kuaidian-app/models"
)

type TaskController struct {
	controllers.BaseController
}

func (c *TaskController) Get() {
	taskId, _ := c.GetInt("taskId", 0)
	task, _ := models.GetTaskById(taskId)
	c.SetJson(0, task, "")
	return

}
