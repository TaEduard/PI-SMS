package main

import (
	"fmt"
	"time"

	"github.com/taeduard/PI-SMS/rabbitmq_rpc"
	"github.com/taeduard/PI-SMS/utils"
)

func main() {
	forever := make(chan bool)
	go func() {
		for {
			res, err := rabbitmq_rpc.RPCcommand("ls")
			time.Sleep(3 * time.Second)
			utils.FailOnError(err, "Failed to handle RPC request")
			fmt.Println(res)
		}
	}()
	<-forever

}
