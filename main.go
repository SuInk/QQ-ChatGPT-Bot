package main

import (
	"QQ-ChatGPT-Bot/cmd/cqhttp"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
func main() {
	go cqhttp.Run()
	for {
	}

}
