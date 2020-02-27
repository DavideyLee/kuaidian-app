package server

import (
	"github.com/astaxie/beego/logs"
	"net/http"

	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/julienschmidt/httprouter"
	//"github.com/xtfly/gokits"
	"io/ioutil"
	"kuaidian-app/library/p2p/common"
	"kuaidian-app/library/p2p/p2p"
	"strconv"
	"strings"
)

func (svc *Server) String(r int, s string, w http.ResponseWriter) {
	w.WriteHeader(r)
	w.Write([]byte(s))
}

func (svc *Server) Json(r int, s interface{}, w http.ResponseWriter) {
	w.WriteHeader(r)
	ss, _ := json.Marshal(s)
	w.Write(ss)
}

func (svc *Server) getRequestParams(r *http.Request, s interface{}) error {
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

// CreateTask POST /api/v1/server/tasks
func (svc *Server) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	t := new(CreateTask)
	if err := svc.getRequestParams(r, t); err != nil {
		log.Errorf("Recv [%s] request, decode body failed. %v", "/api/v1/server/tasks", err)
		return
	}

	// 检查任务是否存在
	v, ok := svc.cache.Get(t.ID)
	if ok {
		cti := v.(*CachedTaskInfo)
		if cti.EqualCmp(t) {
			svc.String(http.StatusAccepted, "", w)
			return
		}
		log.Debugf("[%s] Recv task, task is existed", t.ID)
		svc.String(http.StatusBadRequest, TaskExist.String(), w)
		return
	}

	log.Infof("[%s] Recv task, file=%v, ips=%v", t.ID, t.DispatchFiles, t.DestIPs)

	cti := NewCachedTaskInfo(svc, t)
	svc.cache.Set(t.ID, cti, common.NoExpiration)
	svc.cache.OnEvicted(func(id string, v interface{}) {
		log.Infof("[%s] Remove task cache", t.ID)
		cti := v.(*CachedTaskInfo)
		cti.quitChan <- struct{}{}
	})
	go cti.Start()

	svc.String(http.StatusAccepted, "", w)
	return
}

// CancelTask DELETE /api/v1/server/tasks/:id
func (svc *Server) CancelTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Infof("[%s] Recv cancel task", id)
	v, ok := svc.cache.Get(id)
	if !ok {
		svc.String(http.StatusBadRequest, TaskNotExist.String(), w)
		return
	}
	cti := v.(*CachedTaskInfo)
	cti.stopChan <- struct{}{}
	svc.Json(http.StatusAccepted, "", w)
	return
}

// QueryTask GET /api/v1/server/tasks/:id
func (svc *Server) QueryTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Infof("[%s] Recv query task", id)
	v, ok := svc.cache.Get(id)
	if !ok {
		svc.String(http.StatusBadRequest, TaskNotExist.String(), w)
		return
	}
	cti := v.(*CachedTaskInfo)
	svc.Json(http.StatusOK, cti.Query(), w)
	return
}

// ReportTask POST /api/v1/server/tasks/status
func (svc *Server) ReportTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	csr := new(p2p.StatusReport)
	if err := svc.getRequestParams(r, csr); err != nil {
		log.Errorf("Recv [%s] request, decode body failed. %v", "", err)
		return
	}

	log.Debugf("[%s] Recv task report, ip=%v, pecent=%v", csr.TaskID, csr.IP, csr.PercentComplete)
	if v, ok := svc.cache.Get(csr.TaskID); ok {
		cti := v.(*CachedTaskInfo)
		cti.reportChan <- csr
	}

	svc.String(http.StatusOK, "", w)
	return
}

func (svc *Server) QueryTaskNoHttp(id string) (*TaskInfo, error) {
	log.Infof("[%s] Recv query task", id)
	v, ok := svc.cache.Get(id)
	if !ok {
		return new(TaskInfo), errors.New(TaskNotExist.String())
	}
	cti := v.(*CachedTaskInfo)
	return cti.Query(), nil
}

func (svc *Server) CreateTaskNoHttp(t *CreateTask) error {
	// 检查任务是否存在
	v, ok := svc.cache.Get(t.ID)
	if ok {
		cti := v.(*CachedTaskInfo)
		if cti.EqualCmp(t) {
			return nil
		}
		log.Debugf("[%s] Recv task, task is existed", t.ID)

		return errors.New(TaskExist.String())
	}

	log.Infof("[%s] Recv task, file=%v, ips=%v", t.ID, t.DispatchFiles, t.DestIPs)

	cti := NewCachedTaskInfo(svc, t)
	svc.cache.Set(t.ID, cti, common.NoExpiration)
	svc.cache.OnEvicted(func(id string, v interface{}) {
		log.Infof("[%s] Remove task cache", t.ID)
		cti := v.(*CachedTaskInfo)
		cti.quitChan <- struct{}{}
	})
	go cti.Start()
	return nil
}

// 给所有客户端发送停止命令
func (svc *Server) CheckAllClientIp(ips []string) {
	url := "/api/v1/agent/ip/"
	for _, ip := range ips {
		go func(ip string) {
			sendip := ip
			if idx := strings.Index(ip, ":"); idx > 0 {
				ip = ip[:idx]
			}
			ip = ip + ":" + strconv.Itoa(svc.Cfg.Net.AgentMgntPort)
			if _, err2 := svc.HTTPGet(ip, url+sendip); err2 != nil {
				log.Errorf("Send http request failed. GET, ip=%s, url=%s, error=%v", ip, url, err2)
			} else {
				log.Debugf("Send http request success. GET, ip=%s, url=%s", ip, url)
			}
		}(ip)
	}
}

// 给所有客户端是否存在
func (svc *Server) CheckAllClient(hosts []string) map[string]string {
	res := map[string]string{}
	url := "/api/v1/agent/alive"
	for _, host := range hosts {
		ip := ""
		if idx := strings.Index(host, ":"); idx > 0 {
			ip = host[:idx]
		}
		ip = ip + ":" + strconv.Itoa(svc.Cfg.Net.AgentMgntPort)
		if _, err2 := svc.HTTPGet(ip, url); err2 != nil {
			res[ip] = "dead"
			logs.Error(err2)
		} else {
			res[ip] = "alive"
		}
	}
	return res
}
