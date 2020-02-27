package server

import (
	"github.com/julienschmidt/httprouter"
	"kuaidian-app/library/p2p/common"
	"kuaidian-app/library/p2p/p2p"
	"time"
)

type Server struct {
	common.BaseService
	// 用于缓存当前接收到任务
	cache *common.Cache
	// Session管理
	sessionMgnt *p2p.TaskSessionMgnt
}

func NewServer(cfg *common.Config) (*Server, error) {
	s := &Server{
		cache:       common.NewCache(5 * time.Minute),
		sessionMgnt: p2p.NewSessionMgnt(cfg),
	}
	s.BaseService = *common.NewBaseService(cfg, cfg.Name, s)
	return s, nil
}

func (svc *Server) OnStart(cfg *common.Config, e *httprouter.Router) error {
	go func() { svc.sessionMgnt.Start() }()
	e.POST("/api/v1/server/tasks", svc.CreateTask)
	e.DELETE("/api/v1/server/tasks/:id", svc.CancelTask)
	e.GET("/api/v1/server/tasks/:id", svc.QueryTask)
	e.POST("/api/v1/server/tasks/status", svc.ReportTask)
	return nil
}

func (svc *Server) OnStop(c *common.Config, e *httprouter.Router) {
	go func() { svc.sessionMgnt.Stop() }()
}
