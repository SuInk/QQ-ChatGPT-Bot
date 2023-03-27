package main

import (
	"log"

	"QQ-ChatGPT-Bot/cmd/cqhttp"
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
