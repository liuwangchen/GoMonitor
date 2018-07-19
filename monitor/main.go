package main

import (
	"GoMonitor/monitor/Config"
	"GoMonitor/monitor/Info"
	"GoMonitor/monitor/Operation"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
func WsHandler(w http.ResponseWriter, r *http.Request) (*websocket.Conn, bool) {
	var conn *websocket.Conn
	var err error
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("连接出错：", err)
		return nil, false
	} else {
		address := r.RemoteAddr
		fmt.Println("连上了，地址：", address)
		var connMap map[string]*websocket.Conn
		switch r.RequestURI {
		case "/monitorCpu":
			connMap = connCpuMap
		case "/monitorNet":
			connMap = connNetMap
		case "/monitorProcess":
			connMap = connProcessMap
		}
		_, ok := connMap[address]
		if !ok {
			connMap[address] = conn
		}
		return conn, true
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go runMonitorProcessTicker()
	go runMonitorCpuTicker()
	go runMonitorNetTicker()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	//启动静态服务
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	r.Static("/resources", dir+"/Views/")
	r.LoadHTMLFiles(dir + "/Views/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	//监控cpu
	r.GET("/monitorCpu", func(c *gin.Context) {
		if conn, isConn := WsHandler(c.Writer, c.Request); isConn {
			fmt.Println("当前cpu连接总数：", len(connCpuMap))
			for {
				_, reply, err := conn.ReadMessage()
				if err != nil {
					break
				} else {
					ReceiveCpu(string(reply), conn)
				}
			}
		}
	})
	//监控网络
	r.GET("/monitorNet", func(c *gin.Context) {
		if conn, isConn := WsHandler(c.Writer, c.Request); isConn {
			fmt.Println("当前net连接总数：", len(connNetMap))
			for {
				_, reply, err := conn.ReadMessage()
				if err != nil {
					break
				} else {
					ReceiveNet(string(reply), conn)
				}
			}
		}
	})
	//监控进程
	r.GET("/monitorProcess", func(c *gin.Context) {
		if conn, isConn := WsHandler(c.Writer, c.Request); isConn {
			fmt.Println("当前process连接总数：", len(connProcessMap))
			for {
				_, reply, err := conn.ReadMessage()
				if err != nil {
					break
				} else {
					ReceiveProcess(string(reply), conn)
				}
			}
		}
	})
	r.Run()
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
			cpuInfo := Info.GetCpuInfo()
			//推送
			for k, conn := range connCpuMap {
				go func(k string, conn *websocket.Conn) {
					cpuData := Operation.CpuData(cpuInfo).Monitor(k).Search(k).Sort(k)
					err := conn.WriteJSON(gin.H{
						"total": len(cpuData),
						"rows":  cpuData,
					})
					if err != nil {
						delete(connCpuMap, k)
						delete(Config.OperaCpuConfig, k)
						fmt.Println("当前订阅cpu的连接总数：", len(connCpuMap))
						fmt.Println(k, "cpu用户已断开")
					}
				}(k, conn)
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
			netInfo := Info.GetNetInfo()
			//推送
			for k, conn := range connNetMap {
				go func(k string, conn *websocket.Conn) {
					netData := Operation.NetData(netInfo).Monitor(k).Search(k).Sort(k)
					err := conn.WriteJSON(gin.H{
						"total": len(netData),
						"rows":  netData,
					})
					if err != nil {
						delete(connNetMap, k)
						delete(Config.OperaNetConfig, k)
						fmt.Println("当前订阅net的连接总数：", len(connNetMap))
						fmt.Println(k, "net用户已断开")
					}
				}(k, conn)
			}
		}
	}
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
			processInfo := Info.GetProcessInfo()
			//推送
			for k, conn := range connProcessMap {
				go func(k string, conn *websocket.Conn) {
					processData := Operation.ProcessData(processInfo).Monitor(k).Search(k).Sort(k)
					err := conn.WriteJSON(gin.H{
						"total": len(processData),
						"rows":  processData,
					})
					if err != nil {
						delete(connProcessMap, k)
						delete(Config.OperaProcessConfig, k)
						fmt.Println("当前订阅process的连接总数：", len(connProcessMap))
						fmt.Println(k, "process用户已断开")
					}
				}(k, conn)
			}
		}
	}
}

func ReceiveCpu(message string, conn *websocket.Conn) {
	u, _ := url.Parse(message)
	switch u.Path {
	case "sort":
		{
			property := u.Query().Get("property")
			Config.SetCpuSortConfig(conn.RemoteAddr().String(), property)
		}
	}
}

func ReceiveNet(message string, conn *websocket.Conn) {
	u, _ := url.Parse(message)
	switch u.Path {
	case "sort":
		{
			property := u.Query().Get("property")
			Config.SetNetSortConfig(conn.RemoteAddr().String(), property)
		}
	}
}

func ReceiveProcess(message string, conn *websocket.Conn) {
	u, _ := url.Parse(message)
	switch u.Path {
	case "sort":
		{
			property := u.Query().Get("property")
			Config.SetProcessSortConfig(conn.RemoteAddr().String(), property)
		}
	}
}
