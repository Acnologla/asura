package telemetry

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	masterMetric  = `{ "series" : [%s] }`
	defaultMetric = `{"metric":"asura.%s", "points":[[%d, %d]], "type":"gauge", "host":"asura"}`
	metricURL     = "https://app.datadoghq.com/api/v1/series?api_key=%s"
)

// The main purpose of this function it to send the metrics of the bot to the service "Datadog"
func MetricUpdate() {
	if os.Getenv("production") != "" {
		url := fmt.Sprintf(metricURL, os.Getenv("DATADOG_API_KEY"))
		users := fmt.Sprintf(defaultMetric, "client.users", time.Now().Unix(), 159)
		guilds := fmt.Sprintf(defaultMetric, "client.guilds", time.Now().Unix(), 159)
		channels := fmt.Sprintf(defaultMetric, "client.channels", time.Now().Unix(), 159)
		ram := fmt.Sprintf(defaultMetric, "memory.rss", time.Now().Unix(), 159)
		series := fmt.Sprintf("%s,%s,%s,%s", users, guilds, channels, ram)
		realMetric := fmt.Sprintf(masterMetric, series)
		_, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(realMetric)))
		if err != nil {
			fmt.Println(err)
		}
	}
}
