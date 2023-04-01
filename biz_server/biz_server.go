package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"net/http"
	"os"
	"path"
)

// 帮我们完成升级协议的代理
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool { // 告诉它也是可以跨域
		return true
	},
}

func main() {
	fmt.Println("start bizServer")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	log.Config(path.Dir(ex) + "/log/biz_server.log")
	log.Info("北纬8° log ")

	//http://127.0.0.1:12345/websocket // ab -n 1000 -c 8 http://127.0.0.1:12345/websocket
	//===》 1000 次  8 个并发测   这个只能看每秒的请求次数 看不出CPU 的 负载
	//http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = w.Write([]byte("Hello, the world"))
	//})

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe("127.0.0.1:12345", nil)
}

// websocket 每个请求连接上来了之后都会升级   作用域是队于每个客户端请求(不是每个请求)都是独立的
// 要是我们的客户端 断开了  那么就要把他关闭了
func webSocketHandshake(w http.ResponseWriter, r *http.Request) {
	if w == nil ||
		nil == r {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if nil != err {
		log.Error("Websocket upgrade error,%v+", err)
		return
	}
	//defer conn.Close() //  defer  这个 函数执行到最后的时候执行这个 defer    //  但是我们需要接受这个函数
	defer func() {
		_ = conn.Close() // 直接将异常吃掉
	}()
	log.Info("有新客户端连入")

	for { // 特别简单的就能处理多个连接的情况 他是在不同携程里做的  go 天生支持协程
		_, msgData, err := conn.ReadMessage()
		if nil != err {
			log.Error("%v+", err)
			break
		}

		log.Info("%v", msgData)
		// 解包 消息需要转成我们的消息对象  probuf 相当于是一个公共文档

		cmd := &msg.UserLoginCmd{}
		_ = proto.Unmarshal(msgData[4:], cmd)

		log.Info("%v", cmd)

	}
}
