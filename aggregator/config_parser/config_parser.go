package config_parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigStruct struct {
	Freq     int
	Commands []CommandStruct
}

type CommandStruct struct {
	CommandName string
	Command     string
}

func Parse_config(config string) ConfigStruct {
	jsonFile, err := os.Open(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result ConfigStruct
	json.Unmarshal([]byte(byteValue), &result)

	return result
}
