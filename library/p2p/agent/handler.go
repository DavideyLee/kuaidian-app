package agent

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/julienschmidt/httprouter"
	nettool "github.com/toolkits/net"
	"io/ioutil"
	"kuaidian-app/library/p2p/p2p"
	"net/http"
)

func (c *Agent) getRequestParams(r *http.Request, s interface{}) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	} else {
		if err := json.Unmarshal(rbody, &s); err != nil {
			return err
		}
	}
	return nil

}

//------------------------------------------
// POST /api/v1/agent/tasks
func (c *Agent) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	dt := new(p2p.DispatchTask)
	if err := c.getRequestParams(r, dt); err != nil {
		log.Errorf("Recv '%s' request, decode body failed. %v", "/api/v1/agent/tasks", err)
		return
	}

	log.Infof("[%s] Recv create task request", dt.TaskID)
	// 暂不检查任务是否重复下发
	c.sessionMgnt.CreateTask(dt)
	w.Write([]byte("ok"))
	return
}

// StartTask POST /api/v1/agent/tasks/start
func (c *Agent) StartTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	st := new(p2p.StartTask)
	if err := c.getRequestParams(r, st); err != nil {
		log.Errorf("Recv '%s' request, decode body failed. %v", "/api/v1/agent/tasks/start", err)
		return
	}

	log.Infof("[%s] Recv start task request", st.TaskID)
	// 暂不检查任务是否重复下发
	c.sessionMgnt.StartTask(st)
	w.Write([]byte("ok"))
	return
}

// CancelTask DELETE /api/v1/agent/tasks/:id
func (c *Agent) CancelTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Infof("[%s] Recv cancel task request", id)
	c.sessionMgnt.StopTask(id)
	w.Write([]byte("ok"))
	return
}

// 为了保证本地获取和客户端配置ip的一致性 GET /api/v1/agent/ip/:ip
func (c *Agent) ChangeIp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	ip := ps.ByName("ip")
	if ip != "" {
		LocalIps, err := nettool.IntranetIP()
		if err != nil {
			log.Infof("get LocalIp error")
			return
		}
		for _, LocalIp := range LocalIps {
			if LocalIp == ip {
				c.Cfg.Net.IP = LocalIp
			}
		}
	}
	w.Write([]byte("ok"))
	return
}
func (c *Agent) Alive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("ok"))
	return
}
