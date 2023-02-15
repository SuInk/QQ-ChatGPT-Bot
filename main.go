package main

import (
	"QQ-ChatGPT-Bot/cmd/cqhttp"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
			return
		}
	}()
	go cqhttp.Run()
	for {
		time.Sleep(5 * time.Second)
	}

}
