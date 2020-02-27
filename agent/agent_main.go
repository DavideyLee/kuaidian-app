package main

import (
	"kuaidian-app/library/p2p/agent"
	"kuaidian-app/library/p2p/common"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg := common.ReadJson("agent.json")
	ss, err := common.ParserConfig(&cfg)
	if err != nil {
		log.Printf("conf agent error, %s.\n", err.Error())
		os.Exit(4)
	}
	log.Print("config: ", ss)

	svc, err := agent.NewAgent(&cfg)
	if err != nil {
		log.Printf("start agent error, %s.\n", err.Error())
		os.Exit(4)
	}
	log.Print("agent: ", svc)

	if err = svc.Start(); err != nil {
		log.Printf("Start service failed, %s.\n", err.Error())
		os.Exit(4)
	}
	quitChan := listenSigInt1()
	select {
	case <-quitChan:
		log.Printf("got control-C")
		svc.Stop()
	}
}
func listenSigInt1() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
