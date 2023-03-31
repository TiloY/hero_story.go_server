package main

import (
	"fmt"
	"hero_story.go_server/comm/log"
	"os"
	"path"
)

func main() {
	fmt.Println("start bizServer")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	log.Config(path.Dir(ex) + "/log/biz_server.log")
	log.Info("北纬8° log ")
}
