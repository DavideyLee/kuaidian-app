package wallecontrollers

import (
	"kuaidian-app/controllers"
	"kuaidian-app/library/components"
	"kuaidian-app/models"
)

type JenkinsController struct {
	controllers.BaseController
}

func (c *JenkinsController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{})
	g := components.BasJenkins{}
	g.SetBaseComponents(s)
	res, err := g.GetCommitList()
	if err != nil {
		c.SetJson(1, nil, "获取Commit错误—"+err.Error())
		return
	} else {
		c.SetJson(0, res, "")
		return
	}

}
