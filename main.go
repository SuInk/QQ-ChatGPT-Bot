package main

import (
	"log"

	"QQ-ChatGPT-Bot/cmd/cqhttp"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	go cqhttp.Run()
	for {
		cqhttp.TimeOutCheck()
	}
}
