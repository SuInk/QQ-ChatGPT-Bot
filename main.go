package main

import (
	"QQ-ChatGPT-Bot/cmd/cqhttp"
	"log"
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
		cqhttp.TimeOutCheck()

	}

}
