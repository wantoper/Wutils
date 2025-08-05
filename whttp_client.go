package main

import (
	"WUtils/WHttp"
	"fmt"
	"io"
	"strings"
)

func main() {
	//uri, _ := url.ParseRequestURI("http://localhost:8080/hello")
	//fmt.Println(uri)
	//fmt.Println(uri.Host)
	//fmt.Println(uri.Scheme)
	//fmt.Println(uri.Path)
	//fmt.Println(uri.Port())
	//fmt.Println(uri.Query())
	header := WHttp.NewHeader()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("ApiKey", "api-key-12345")

	//response := WHttp.Get("http://localhost:8080/hello", header)
	//if response != nil {
	//	fmt.Println("Response Status Code:", response.StatusCode)
	//}
	//fmt.Println("Response Headers:", response.Header)
	//body, err := io.ReadAll(response.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body:", err)
	//} else {
	//	fmt.Println("Response Body:", string(body))
	//}

	reader := strings.NewReader("{'name':'WHttp','version':'1.0.0'}")
	response := WHttp.Post("http://localhost:8080/api/data", header, reader)
	if response != nil {
		fmt.Println("Response Status Code:", response.StatusCode)
	}
	fmt.Println("Response Headers:", response.Header)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("Response Body:", string(body))
	}

	//client := http.DefaultClient
	//request, err := http.NewRequest("GET", "http://localhost:8080/hello", nil)
	//client.Get()
	//client.Do(request)
}
