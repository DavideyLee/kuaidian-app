package init_sever

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"kuaidian-app/library/p2p/common"
	"kuaidian-app/library/p2p/server"
	"os"
)

var P2pSvc *server.Server

func init() {

}

func Start() {
	cfg := common.ReadJson("agent/server.json")
	_, err := common.ParserConfig(&cfg)
	cfg.Server = true
	P2pSvc, err = server.NewServer(&cfg)
	if err != nil {
		logs.Error("start server error, %s.\n", err.Error())
		if beego.BConfig.RunMode != "docker" {
			os.Exit(4)
		}
	}
	logs.Info("服务端p2p配置检测成功")
	if err := P2pSvc.Start(); err != nil {
		logs.Error("Start service failed, %s.\n", err.Error())
		if beego.BConfig.RunMode != "docker" {
			os.Exit(4)
		}
	}
}
