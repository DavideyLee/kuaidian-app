package p2pcontrollers

import (
	"github.com/astaxie/beego/logs"
	"kuaidian-app/controllers"
	"kuaidian-app/library/components"
	"kuaidian-app/models"

	"github.com/astaxie/beego"
)

type SendAgentController struct {
	controllers.BaseController
}

func (c *SendAgentController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}

	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -3})
	agentDir := beego.AppConfig.String("AgentDir")
	AgentDestDir := beego.AppConfig.String("AgentDestDir")

	err := s.SendP2pAgent(agentDir, AgentDestDir)
	logs.Warning(agentDir, AgentDestDir)
	if err != nil {
		logs.Error("出错啦！")
		logs.Error(err.Error())
		c.SetJson(1, nil, "p2p文件传输失败，请检查配置，或目标机器权限"+err.Error())
		return
	}
	c.SetJson(0, nil, "更新agent成功")
	return
}
