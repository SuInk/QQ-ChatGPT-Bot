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
	go cqhttp.Run()
	for {
		time.Sleep(5 * time.Second)
	}

}
