package wallecontrollers

import (
	"kuaidian-app/controllers"
	"kuaidian-app/library/common"
	"kuaidian-app/library/components"
	"kuaidian-app/models"
	"strings"
)

type FlushController struct {
	controllers.BaseController
}

func (c *FlushController) Get() {
	projectIds := c.GetString("projectIds")
	projectIdsArr := strings.Split(projectIds, ",")
	res := []map[string]interface{}{}
	for _, projectId := range projectIdsArr {
		Project, err := models.GetProjectById(common.GetInt(projectId))
		if err != nil {
			continue
		}
		s := components.BaseComponents{}
		s.SetProject(Project)
		s.SetTask(&models.Task{Id: -2})
		err = s.GetExecFlush()
		if err != nil {
			res = append(res, map[string]interface{}{"name": Project.Name, "err": err.Error()})
		} else {
			res = append(res, map[string]interface{}{"name": Project.Name, "msg": "success"})
		}
	}
	c.SetJson(0, res, "")
	return

}
