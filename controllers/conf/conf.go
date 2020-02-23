package confcontrollers

import (
	"github.com/astaxie/beego/logs"
	"kuaidian-app/controllers"

	"encoding/json"
	"kuaidian-app/models"
)

type ConfController struct {
	controllers.BaseController
}

func (c *ConfController) Get() {
	projectId, _ := c.GetInt("projectId", 0)
	project, _ := models.GetProjectById(projectId)
	c.SetJson(0, project, "")
	return

}
func (c *ConfController) Post() {
	//projectId,_:=c.GetInt("projectId",0)
	var project models.Project
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &project)
	err = models.UpdateProjectById(&project)
	logs.Info(err)
	return
}
