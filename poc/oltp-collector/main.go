package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

var meterProvider *metric.MeterProvider

func init() {

}

func main() {
	ctx := context.Background()
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure())
	if err != nil {
		panic(err)
	}

	meterProvider = metric.NewMeterProvider(metric.WithReader(
		metric.NewPeriodicReader(exporter)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	global.SetMeterProvider(meterProvider)

	r := gin.Default()

	commonLabels := []attribute.KeyValue{
		attribute.String("test-attr", "test-value"),
	}

	r.GET("/hello", func(c *gin.Context) {
		meter := global.Meter("oltp-hello-metrics")
		requestCount, _ := meter.SyncFloat64().Counter("oltp-hello-metrics/incoming-request", instrument.WithDescription("request processed"))
		requestCount.Add(c, 1, commonLabels...)
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	r.Run()
}
