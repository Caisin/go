package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	APIURL = "https://free-api.heweather.net/s6/air/now"
	//本地
	APIKEY = "a5665e0243a345e8a57d46dca0230560"
	//生产
	//APIKEY = "a6949c2baba743788cddc8a9f5960cbb"
)

func main() {
	queryUrl := fmt.Sprintf("%s?key=%s&amp;location=beijing", APIURL, APIKEY)
	resp, err := http.Get(queryUrl)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
