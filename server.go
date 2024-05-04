package main

import (
	"encoding/json"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func NewSumTaskSignature(val []string, sessionID string) *tasks.Signature {
	// 新建任务签名，执行任务名Name，指定传递给任务的参数Args
	signature := &tasks.Signature{
		Name: "SayHello",
		Args: []tasks.Arg{
			{
				Type:  "[]string",
				Value: val,
			},
			{
				Type:  "string",
				Value: sessionID,
			},
		},
	}
	return signature
}

// websocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var server *machinery.Server // 全局变量，用于存储服务器实例

var err error

type Receive struct {
	Mode string   `json:"mode"`
	Arg  []string `json:"arg"`
}

func init() {
	// 加载配置文件
	cnf, err := config.NewFromYaml("./config.yaml", false)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err) // 使用 log.Fatalf 来记录错误并停止程序
	}

	// 创建服务器实例
	server, err = machinery.NewServer(cnf)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
}

func main() {
	//启动worker
	go mainWorker()
	
	// 创建HTTP服务器
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("服务器开启在端口:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 定义一个全局的连接映射和锁
//var connections = make(map[string]*websocket.Conn)
//var lock sync.RWMutex

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, errGrader := upgrader.Upgrade(w, r, nil)
	if errGrader != nil {
		log.Println("连接升级失败：", errGrader)
		return
	}

	// 生成并注册连接
	sessionID := generateSessionID()

	// 注册连接
	lock.Lock()
	connections[sessionID] = conn
	lock.Unlock()

	// 确保在结束时注销连接
	defer func() {
		lock.Lock()
		delete(connections, sessionID)
		lock.Unlock()
		if err := conn.Close(); err != nil {
			log.Printf("关闭连接失败: %v", err)
		}
	}()

	for {
		messageType, p, errRead := conn.ReadMessage()
		if errRead != nil {
			log.Println("读取消息失败/连接关闭", errRead)
			return
		}
		log.Println("收到消息:", string(p))

		// 将收到消息转换为 JSON
		var receive Receive
		err = json.Unmarshal(p, &receive)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}
		log.Println("收到消息，转为JSON：", receive)

		if receive.Mode == "task" {
			// 下发任务给消费者
			signature := NewSumTaskSignature(receive.Arg, sessionID)
			_, err := server.SendTask(signature)
			if err != nil {
				log.Fatal("下发任务失败：", err)
			}

			// 每秒获取一次消息队列中的结果
			//res, err := asyncResult.Get(1)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//log.Printf("get res is %v\n", tasks.HumanReadableResults(res))
		}

		errSend := conn.WriteMessage(messageType, []byte("Hello!"))
		if errSend != nil {
			log.Println("发送数据失败：", errSend)
			return
		}
	}
}

func generateSessionID() string {
	return uuid.New().String()
}
