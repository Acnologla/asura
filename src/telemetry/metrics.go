package telemetry

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/andersfylling/disgord"
)

const (
	masterMetric  = `{ "series" : [%s] }`
	defaultMetric = `{"metric":"asura.%s", "points":[[%d, %d]], "type":"gauge", "host":"asura"}`
	metricURL     = "https://app.datadoghq.com/api/v1/series?api_key=%s"
)

// The main purpose of this function it to send the metrics of the bot to the service "Datadog"
func metricUpdate(session disgord.Session) {
	if os.Getenv("PRODUCTION") != "" {
		url := fmt.Sprintf(metricURL, os.Getenv("DATADOG_API_KEY"))
		date := time.Now().Unix()
		guildsSize := session.GetConnectedGuilds()
		var memory runtime.MemStats
		runtime.ReadMemStats(&memory)
		guilds := fmt.Sprintf(defaultMetric, "client.guilds", date, len(guildsSize))
		ram := fmt.Sprintf(defaultMetric, "memory.rss", date, memory.Alloc/1000/1000)
		ping, _ := session.HeartbeatLatencies()
		realPing := fmt.Sprintf(defaultMetric, "client.ping", date, ping[0].Milliseconds())
		series := fmt.Sprintf("%s,%s,%s", guilds, ram, realPing)
		realMetric := fmt.Sprintf(masterMetric, series)
		res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(realMetric)))
		if err != nil {
			Error(err.Error(), map[string]string{})
			return
		}
		res.Body.Close()
	}
}

func MetricUpdate(session disgord.Session) {
	for {
		metricUpdate(session)
		time.Sleep(5 * time.Minute)
	}
}
