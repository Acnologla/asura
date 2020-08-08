package main

import (
	"net/http"
	"fmt"
    "io/ioutil"
)

func main(){
	req, _ := http.NewRequest("HEAD", "https://golang.org/pkg/net/http/#Header", nil)
    fmt.Println(req)
    var client http.Client
    resp, _ := client.Do(req)
    fmt.Println(resp.Header["Content-Type"])
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(len(body))
}
