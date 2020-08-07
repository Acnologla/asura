package logs

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"bytes"
)

var url string

func log(level string, message string,values map[string]string){
	values["message"] = message
	values["level"] = level
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[%s] - %s\n",level,message)
}

func Debug(message string,values map[string]string){
	log("debug",message,values)
}


func Info(message string, values map[string]string){
	log("info",message,values)
}

func Warn(message string,values map[string]string){
	log("warn",message,values)
}

func Error(message string,values map[string]string){
	log("error",message,values)
}

func Init(){
	url = fmt.Sprintf("https://http-intake.logs.datadoghq.com/v1/input/%s?ddsource=nodejs&service=asura",os.Getenv("DATADOG_API_KEY"))
}