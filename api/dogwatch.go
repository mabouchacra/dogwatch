package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type SimpleMetric struct {
	name   string
	volume int64
}

func initDDcnx() (context.Context, *datadog.APIClient) {
	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: os.Getenv("DD_API_KEY"),
			},
			"appKeyAuth": {
				Key: os.Getenv("DD_APPLICATION_KEY"),
			},
		},
	)

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	return ctx, apiClient
}

func GetLogBasedMetricVolume() {
	GetLogBasedMetricVolumeWithLimit(-1)
}

func GetLogBasedMetricVolumeWithLimit(lowerLimit int) []SimpleMetric {
	ctx, apiClient := initDDcnx()

	metricsResponse := getLogBasedMetrics(ctx, apiClient)

	var metrics []SimpleMetric

	for _, metric := range metricsResponse {
		response := getMetricVolume(metric.GetId(), ctx, apiClient)
		if response.Data.MetricDistinctVolume != nil {

			volume := *response.Data.MetricDistinctVolume.Attributes.DistinctVolume
			if volume >= int64(lowerLimit) {
				metrics = append(metrics, SimpleMetric{metric.GetId(), volume})
			}
		} else {
			responseContent, _ := json.MarshalIndent(response, "", "  ")
			fmt.Fprintf(os.Stderr, "Error deserializing metric %s - raw value : %s\n", metric.GetId(), &responseContent)
		}
	}
	return metrics
}

func getLogBasedMetrics(ctx context.Context, apiClient *datadog.APIClient) []datadogV2.LogsMetricResponseData {
	api := datadogV2.NewLogsMetricsApi(apiClient)
	response, r, err := api.ListLogsMetrics(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling 'LogsMetricsApi.ListLogsMetrics' : %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response : %v\n", r)
	}
	return response.Data
}

func getMetricVolume(metric string, ctx context.Context, apiClient *datadog.APIClient) datadogV2.MetricVolumesResponse {
	api := datadogV2.NewMetricsApi(apiClient)
	resp, r, err := api.ListVolumesByMetricName(ctx, metric)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.ListVolumesByMetricName`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	return resp
}
