package p2pcontrollers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"kuaidian-app/controllers"
	"kuaidian-app/library/components"
	"kuaidian-app/library/p2p/init_sever"
	"kuaidian-app/models"
)

type AgentController struct {
	controllers.BaseController
}

func (c *AgentController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -3})
	Hosts := s.GetAllHost()
	ss := init_sever.P2pSvc.CheckAllClient(Hosts)
	reHosts := []string{}
	for host, status := range ss {
		if status == "dead" {
			reHosts = append(reHosts, host)
		}
	}
	if len(reHosts) > 0 && c.Project.P2p == 1 {
		AgentDestDir := beego.AppConfig.String("AgentDestDir")
		err := s.StartP2pAgent(reHosts, AgentDestDir)
		if err != nil {
			c.SetJson(1, nil, "重启失败"+err.Error())
			return
		} else {
			c.SetJson(0, nil, "重启成功")
			return
		}
	} else {
		c.SetJson(0, nil, "已全部启动")
		return
	}
}

func (c *AgentController) Post() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	ips := []string{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &ips)
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	AgentDestDir := beego.AppConfig.String("AgentDestDir")
	err := s.StartP2pAgent(ips, AgentDestDir)
	if err != nil {
		c.SetJson(1, nil, "重启失败"+err.Error())
		return
	} else {
		c.SetJson(0, nil, "重启成功")
		return
	}
}
