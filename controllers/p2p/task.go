package p2pcontrollers

import (
	"kuaidian-app/controllers"
	"kuaidian-app/library/p2p/init_sever"
)

type TaskController struct {
	controllers.BaseController
}

func (c *TaskController) Get() {
	taskId := c.GetString("taskId")
	ss, _ := init_sever.P2pSvc.QueryTaskNoHttp(taskId)
	c.SetJson(0, ss, "")
	return
}
