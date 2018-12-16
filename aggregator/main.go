package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/taeduard/PI-SMS/aggregator/config_parser"
	"github.com/taeduard/PI-SMS/rabbitmq_rpc"
	"github.com/taeduard/PI-SMS/utils"
)

// map[hostname][command][timestamp]result
type ResultStruct map[string]map[string]map[string]string

func main() {
	cnfig := config_parser.Parse_config("config.json")
	resStruct := make(ResultStruct)
	forever := make(chan bool)
	file, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	for _, elem := range cnfig.Commands {
		go func(command string, commandName string) {
			for {
				// fmt.Println("command:", command)
				// fmt.Println("commandName:", commandName)
				res, err := rabbitmq_rpc.RPCcommand(command)
				time.Sleep(time.Duration(cnfig.Freq) * time.Second)
				utils.FailOnError(err, "Failed to handle RPC request")
				b := strings.Split(res, ";")
				resStruct[b[1]] = make(map[string]map[string]string)
				resStruct[b[1]][command] = make(map[string]string)
				resStruct[b[1]][command][b[0]] = b[2]
				out, err := json.Marshal(resStruct)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(file, fmt.Sprintf("\n %s\n %s \n", commandName, string(out)))
			}
		}(elem.Command, elem.CommandName)
	}

	<-forever

}
