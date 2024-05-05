package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func mainCron(args []string, sessionID string) {
	// 使用秒级别的定时器
	c := cron.New(cron.WithSeconds())

	//// 添加任务，每分钟的第0秒执行
	//_, err := c.AddFunc("0 * * * * *", func() {
	//	// 下发任务给消费者
	//	signature := NewSumTaskSignature(args, sessionID)
	//	_, err := server.SendTask(signature)
	//	if err != nil {
	//		log.Fatal("下发任务失败：", err)
	//	}
	//})
	//if err != nil {
	//	fmt.Println("添加定时任务出错:", err)
	//	return
	//}

	// 添加任务, 每天的15:30执行
	_, err := c.AddFunc("0 30 15 * * *", func() {
		// 下发任务给消费者
		signature := NewSumTaskSignature(args, sessionID)
		_, err := server.SendTask(signature)
		if err != nil {
			log.Fatal("下发任务失败：", err)
		}
		fmt.Println("任务执行中:", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Println("添加定时任务出错:", err)
		return
	}

	// 启动cron调度器
	c.Start()

	// 保持主程序运行，以便任务可以按计划执行
	select {}
}
