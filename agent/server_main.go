package main

import (
	"kuaidian-app/library/p2p/common"
	"kuaidian-app/library/p2p/server"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg := common.ReadJson("server.json")
	ss, err := common.ParserConfig(&cfg)
	if err != nil {
		log.Printf("conf server error, %s.\n", err.Error())
		os.Exit(4)
	}
	log.Print("config: ", ss)

	svc, err := server.NewServer(&cfg)
	if err != nil {
		log.Printf("start server error, %s.\n", err.Error())
		os.Exit(4)
	}
	log.Print("server: ", svc)

	if err = svc.Start(); err != nil {
		log.Printf("Start service failed, %s.\n", err.Error())
		os.Exit(4)
	}

	quitChan := listenSigInt()
	select {
	case <-quitChan:
		log.Printf("got control-C")
		svc.Stop()
	}
}

func listenSigInt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
