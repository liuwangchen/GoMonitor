package main

import (
	"GoMonitor/Model"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//订阅cpu的conn集合
var connCpuMap = make(map[string]*websocket.Conn)

//订阅net的conn集合
var connNetMap = make(map[string]*websocket.Conn)

//订阅process的conn集合
var connProcessMap = make(map[string]*websocket.Conn)

// 处理ws请求
func WsHandler(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("连接出错：", err)
		return
	} else {
		remoteAddr := conn.RemoteAddr().String()
		fmt.Println("连上了，地址：", remoteAddr)
		var connMap map[string]*websocket.Conn
		switch r.RequestURI {
		case "/monitorCpu":
			connMap = connCpuMap
		case "/monitorNet":
			connMap = connNetMap
		case "/monitorProcess":
			connMap = connProcessMap
		}
		_, ok := connMap[remoteAddr]
		if !ok {
			connMap[remoteAddr] = conn
		}
		fmt.Println("当前cpu连接总数：", len(connCpuMap))
		fmt.Println("当前net连接总数：", len(connNetMap))
		fmt.Println("当前process连接总数：", len(connProcessMap))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go runMonitorProcessTicker()
	go runMonitorCpuTicker()
	go runMonitorNetTicker()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Static("/resources","../Views/")
	r.LoadHTMLGlob("../Views/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/monitorCpu", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})
	r.GET("/monitorNet", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})
	r.GET("/monitorProcess", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})
	r.Run()
}

//监控进程数据
func runMonitorProcessTicker() {
	defer func() {
		if result := recover(); result != nil {
			log.Println(result)
		}
	}()
	for range time.NewTicker(time.Second * 1).C {
		if len(connProcessMap) > 0 {
			processInfo := Model.GetProcessInfo()
			//推送
			for k, conn := range connProcessMap {
				err := conn.WriteJSON(processInfo)
				if err != nil {
					delete(connProcessMap, k)
					fmt.Println("当前订阅process的连接总数：", len(connProcessMap))
					fmt.Println(conn.RemoteAddr().String(), "process用户已断开")
				}
			}
		}
	}
}

//监控网路数据
func runMonitorNetTicker() {
	defer func() {
		if result := recover(); result != nil {
			log.Println(result)
		}
	}()
	for range time.NewTicker(time.Second * 1).C {
		if len(connNetMap) > 0 {
			netInfo := Model.GetNetInfo()
			//推送
			for k, conn := range connNetMap {
				err := conn.WriteJSON(netInfo)
				if err != nil {
					delete(connNetMap, k)
					fmt.Println("当前订阅net的连接总数：", len(connNetMap))
					fmt.Println(conn.RemoteAddr().String(), "net用户已断开")
				}
			}
		}
	}
}

//监控cpu数据
func runMonitorCpuTicker() {
	defer func() {
		if result := recover(); result != nil {
			log.Println(result)
		}
	}()
	for range time.NewTicker(time.Second * 2).C {
		if len(connCpuMap) > 0 {
			cpuInfo := Model.GetCpuInfo()
			//推送
			for k, conn := range connCpuMap {
				err := conn.WriteJSON(cpuInfo)
				if err != nil {
					delete(connCpuMap, k)
					fmt.Println("当前订阅cpu的连接总数：", len(connCpuMap))
					fmt.Println(conn.RemoteAddr().String(), "cpu用户已断开")
				}
			}
		}
	}
}
