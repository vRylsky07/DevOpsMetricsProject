package main

import (
	"DevOpsMetricsProject/internal/metrics"
	"fmt"
)

func main() {
	test := metrics.GetGaugeMetrics()
	test2 := metrics.GetCounterMetrics()
	test2 = metrics.GetCounterMetrics()
	test2 = metrics.GetCounterMetrics()
	for k, v := range test {
		fmt.Println("Key: "+k+" Value: ", v)
	}
	fmt.Println("POLL COUNT: ", test2["PollCount"])
}
