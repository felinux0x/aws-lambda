package observability

import (
	"encoding/json"
	"fmt"
	"time"
)

// EMF represents the AWS CloudWatch Embedded Metric Format.
// Using this directly avoids heavy dependencies and speeds up cold starts.
type EMF struct {
	AWS struct {
		Timestamp int64             `json:"Timestamp"`
		Metrics   []MetricDirective `json:"CloudWatchMetrics"`
	} `json:"_aws"`
	// Custom dimensions and metrics go at the root level
	Dimensions map[string]string  `json:"-"`
	Metrics    map[string]float64 `json:"-"`
}

type MetricDirective struct {
	Namespace  string             `json:"Namespace"`
	Dimensions [][]string         `json:"Dimensions"`
	Metrics    []MetricDefinition `json:"Metrics"`
}

type MetricDefinition struct {
	Name string `json:"Name"`
	Unit string `json:"Unit"`
}

// LogEMF formats and prints metrics synchronously so CloudWatch can ingest them asynchronously without latency overhead.
func LogEMF(namespace string, dimensions map[string]string, metrics map[string]float64) {
	emf := EMF{}
	emf.AWS.Timestamp = time.Now().UnixMilli()

	// Build Dimensions structure
	var dimKeys []string
	for k := range dimensions {
		dimKeys = append(dimKeys, k)
	}

	// Build Metrics structure
	var metricDefs []MetricDefinition
	for k := range metrics {
		metricDefs = append(metricDefs, MetricDefinition{Name: k, Unit: "Count"})
	}

	emf.AWS.Metrics = []MetricDirective{
		{
			Namespace:  namespace,
			Dimensions: [][]string{dimKeys},
			Metrics:    metricDefs,
		},
	}

	// Flatten dimensions and metrics to root level for JSON marshalling
	payload := make(map[string]interface{})
	payload["_aws"] = emf.AWS
	for k, v := range dimensions {
		payload[k] = v
	}
	for k, v := range metrics {
		payload[k] = v
	}

	bytes, err := json.Marshal(payload)
	if err == nil {
		fmt.Println(string(bytes)) // CloudWatch agent reads stdout to ingest EMF
	}
}
