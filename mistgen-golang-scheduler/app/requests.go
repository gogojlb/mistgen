package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SensorReqbody struct {
	Props []string `json:"props"`
}

type SwitchReqBody struct {
	IsOn bool `json:"is_on"`
}

func NewSensorReq(miApiUrl string) *http.Request {

	b := SensorReqbody{
		Props: []string{"temperature", "humidity"},
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/sensor", miApiUrl)
	bBuffer := bytes.NewBuffer(bBytes)

	req, err := http.NewRequest("GET", url, bBuffer)

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "/")

	return req
}

func NewSwitchReq(miApiUrl string, isOn bool) *http.Request {

	b := SwitchReqBody{IsOn: isOn}
	bBytes, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s/switch", miApiUrl)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bBytes))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req
}
