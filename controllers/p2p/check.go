package p2pcontrollers

import (
	"github.com/astaxie/beego/logs"
	"kuaidian-app/controllers"
	"kuaidian-app/library/common"
	"kuaidian-app/library/components"
	"kuaidian-app/library/p2p/init_sever"
	"kuaidian-app/models"

	"github.com/astaxie/beego/orm"
)

type CheckController struct {
	controllers.BaseController
}

type P2pinfo struct {
	Host   string
	Status string
	Pid    int
	Pname  string
	Eid int16
}

func (c *CheckController) Get() {
	searchtype := c.GetString("type")
	projectName := c.GetString("projectName")
	logs.Info(searchtype)
	if searchtype == "0" {
		o := orm.NewOrm()
		var projects []models.Project
		var p []P2pinfo
		ss := map[string]string{}
		i, err := o.Raw("SELECT * FROM `project` WHERE `p2p` = 1 ").QueryRows(&projects)
		if i > 0 && err == nil {
			for _, project := range projects {
				s := components.BaseComponents{}
				s.SetProject(&project)
				ips := s.GetAllHost()
				proRes := init_sever.P2pSvc.CheckAllClient(ips)
				for key, value := range proRes {
					if value == "dead" {
						pa := P2pinfo{}
						if !common.InList(key, ss) {
							ss[key] = value
							pa.Host = key
							pa.Status = value
							pa.Pid = project.Id
							pa.Pname = project.Name
							pa.Eid = project.Level
							p = append(p, pa)
						}
					}
				}
			}
			c.SetJson(0, p, "")
			return
		} else {
			c.SetJson(1, ss, "no agent")
			return
		}
	} else if projectName != "" && searchtype == "1" {
		o := orm.NewOrm()
		var projects []models.Project
		var p []P2pinfo
		ss := map[string]string{}
		i, err := o.Raw("SELECT * FROM `project` WHERE `p2p` = 1 and `name` = ?   ", projectName).QueryRows(&projects)
		if i > 0 && err == nil {
			for _, project := range projects {
				s := components.BaseComponents{}
				s.SetProject(&project)
				ips := s.GetAllHost()
				proRes := init_sever.P2pSvc.CheckAllClient(ips)
				for key, value := range proRes {
					pa := P2pinfo{}
					if !common.InList(key, ss) {
						ss[key] = value
						pa.Host = key
						pa.Status = value
						pa.Pid = project.Id
						pa.Pname = project.Name
						pa.Eid = project.Level
						p = append(p, pa)
					}
				}
			}
			c.SetJson(0, p, "")
			return
		} else {
			c.SetJson(1, ss, "no agent")
			return
		}
	}
	return
}
