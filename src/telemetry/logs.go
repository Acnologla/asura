package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var url string

// Beautify the log message and send it the DataDog service
func log(level string, message string, values map[string]string) {
	if os.Getenv("PRODUCTION") != "" {
		values["message"] = message
		values["level"] = level
		jsonValue, _ := json.Marshal(values)
		res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
	fmt.Printf("[%s] - %s\n", level, message)
}

// Different "flavours" of log message to make it easier to separate

func Debug(message string, values map[string]string) {
	go log("debug", message, values)
}

func Info(message string, values map[string]string) {
	go log("info", message, values)
}

func Warn(message string, values map[string]string) {
	go log("warn", message, values)
}

func Error(message string, values map[string]string) {
	go log("error", message, values)
}

// Initialize the url variable with the DataDogKey that is found in .env or the environment variables.
func Init() {
	url = fmt.Sprintf("https://http-intake.logs.datadoghq.com/v1/input/%s?ddsource=nodejs&service=asura", os.Getenv("DATADOG_API_KEY"))
}
