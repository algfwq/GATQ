package main

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// SayHello 定义任务
func SayHello(args []string, sessionID string) (string, error) {
	for _, arg := range args {
		//time.Sleep(time.Second * 10)
		t := time.Now()
		fmt.Println(arg, t.Format("02 January 2006 15:04"))
	}
	lock.RLock()
	conn, exists := connections[sessionID]
	lock.RUnlock()

	time.Sleep(time.Second * 10)

	if exists {
		t := time.Now()
		err := conn.WriteMessage(websocket.TextMessage, []byte("OK!!!!!!!!!!"+t.Format("02 January 2006 15:04")))
		if err != nil {
			log.Println("Failed to send message:", err)
		}
	} else {
		log.Println("WebSocket connection does not exist for session ID:", sessionID)
	}

	return "ok", nil
}

func mainWorker() {
	// 将配置文件实例化
	cnf, err := config.NewFromYaml("./config.yaml", false)
	if err != nil {
		log.Println("config.NewFromYaml failed, err:", err)
		return
	}

	// 根据实例化的配置文件创建 server 实例
	server, err := machinery.NewServer(cnf)
	if err != nil {
		log.Println("machinery.NewServer failed, err:", err)
		return
	}

	// 为消费者程序注册任务
	err = server.RegisterTask("SayHello", SayHello)
	if err != nil {
		log.Println("server.RegisterTask failed, err:", err)
		return
	}

	// 创建 worker 实例并绑定任务队列名
	worker := server.NewWorker("sms", 1)
	// 运行 worker 监听逻辑，监听消息队列中的任务
	err = worker.Launch()
	if err != nil {
		log.Println("start worker error", err)
		return
	}
}
