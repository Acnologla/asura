package telemetry

import (
	"bytes"
	"fmt"
	"net/http"
	"github.com/andersfylling/disgord"
	"os"
	"runtime"
	"context"
	"time"
)

const (
	masterMetric  = `{ "series" : [%s] }`
	defaultMetric = `{"metric":"asura.%s", "points":[[%d, %d]], "type":"gauge", "host":"asura"}`
	metricURL     = "https://app.datadoghq.com/api/v1/series?api_key=%s"
)

// The main purpose of this function it to send the metrics of the bot to the service "Datadog"
func metricUpdate(client *disgord.Client) {
	if os.Getenv("PRODUCTION") != ""{
		url := fmt.Sprintf(metricURL, os.Getenv("DATADOG_API_KEY"))
		date := time.Now().Unix()
		guildsSize,err := client.GetGuilds(context.Background(),&disgord.GetCurrentUserGuildsParams{})
		if err != nil {
			Error(err.Error(),map[string]string{})
			return
		}
		var memory runtime.MemStats
        runtime.ReadMemStats(&memory)
		guilds := fmt.Sprintf(defaultMetric, "client.guilds", date, len(guildsSize))
		ram := fmt.Sprintf(defaultMetric, "memory.rss", date,memory.TotalAlloc / 1024 / 1024)
		series := fmt.Sprintf("%s,%s", guilds, ram)
		realMetric := fmt.Sprintf(masterMetric, series)
		res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(realMetric)))
		if err != nil {
			Error(err.Error(),map[string]string{})
		}
		res.Body.Close()
	}
}

func MetricUpdate(client *disgord.Client){
	for {
		metricUpdate(client)
		time.Sleep(5 * time.Minute)
	}
}