package main

import (
	"GoMonitor/Info"
	"GoMonitor/UserSort"
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
	r.Static("/resources", "../Views/")
	r.LoadHTMLFiles("../Views/index.html")
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
				cpuData := UserSort.CpuData(cpuInfo).Sort(k)
				err := conn.WriteJSON(gin.H{
					"total": len(cpuData),
					"rows":  cpuData,
				})
				if err != nil {
					delete(connCpuMap, k)
					delete(UserSort.SortCpuConfig, k)
					fmt.Println("当前订阅cpu的连接总数：", len(connCpuMap))
					fmt.Println(k, "cpu用户已断开")
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
			netInfo := Info.GetNetInfo()
			//推送
			for k, conn := range connNetMap {
				netData := UserSort.NetData(netInfo).Sort(k)
				err := conn.WriteJSON(gin.H{
					"total": len(netData),
					"rows":  netData,
				})
				if err != nil {
					delete(connNetMap, k)
					delete(UserSort.SortNetConfig, k)
					fmt.Println("当前订阅net的连接总数：", len(connNetMap))
					fmt.Println(k, "net用户已断开")
				}
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
				processData := UserSort.ProcessData(processInfo).Sort(k)
				err := conn.WriteJSON(gin.H{
					"total": len(processData),
					"rows":  processData,
				})
				if err != nil {
					delete(connProcessMap, k)
					delete(UserSort.SortProcessConfig, k)
					fmt.Println("当前订阅process的连接总数：", len(connProcessMap))
					fmt.Println(k, "process用户已断开")
				}
			}
		}
	}
}

func ReceiveCpu(message string, conn *websocket.Conn) {
	//UserSort.SetCpuSortConfig(conn.RemoteAddr().String(), c.DefaultQuery("propertyName", ""), c.DefaultQuery("sort", "asc"))
}

func ReceiveNet(message string, conn *websocket.Conn) {
	//UserSort.SetCpuSortConfig(c.DefaultQuery("uuid", ""), c.DefaultQuery("propertyName", ""), c.DefaultQuery("sort", "asc"))
}

func ReceiveProcess(message string, conn *websocket.Conn) {
	//UserSort.SetCpuSortConfig(c.DefaultQuery("uuid", ""), c.DefaultQuery("propertyName", ""), c.DefaultQuery("sort", "asc"))
}
