package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/apex/go-apex"
)

type message struct {
	SerialNumber   string `json:"serialNumber"`
	BatteryVoltage string `json:"batteryVoltage"`
	ClickType      string `json:"clickType"`
}

func main() {
	log.SetOutput(os.Stderr)

	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		log.Println(m)

		return m, nil
	})
}
