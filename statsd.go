package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
)

func giveMetricsToDatadogStatsd(configs Configs, metrics []string, model string) {
	client, err := statsd.New(configs.StatsdIPPort) // udp

	if err != nil {
		logFatalf("Error creating StatsD client: %v", err)
	}

	defer client.Close()

	retrievedFloatMetrics := map[string]float64{}

	for _, metric := range metrics {
		metric = strings.ToLower(strings.TrimSpace(metric))
		retrievedMetric := strings.Split(metric, "=")[0]
		retrievedMetricValue := strings.Split(metric, "=")[1]

		if strings.Contains(retrievedMetricValue, ".0") {
			retrievedMetricFloat, err := strconv.ParseFloat(retrievedMetricValue, 64)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			retrievedFloatMetrics[retrievedMetric] = retrievedMetricFloat
		}
	}

	for key, value := range retrievedFloatMetrics {
		err = client.Gauge(key, value, []string{"gateway:" + model}, 1)
		if err != nil {
			log.Printf("Error sending %s to statsd: %v", key, err)
		}
	}
}
