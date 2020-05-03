package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestWeather(t *testing.T) {

	//城市列表 https://a.hecdn.net/download/dev/china-city-list.csv
	response, err := http.Get("http://t.weather.sojson.com/api/weather/city/101250109")
	if err != nil {
		log.Print(err.Error())
	}
	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	var netReturn map[string]interface{}
	err = json.Unmarshal(all, &netReturn)
	println(netReturn)
	println(string(all[:]))
}
