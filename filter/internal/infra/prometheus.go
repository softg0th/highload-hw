package infra

import (
	"context"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	SpamMessagesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "spam_messages_total",
			Help: "Total number of the spam messages",
		})
	SpamBatchesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "spam_batches_total",
			Help: "Total number of the spam batches",
		})
)

func StartPrometheus(
	ctx context.Context,
	prometheusPort string,
	logger *logstash_logger.Logstash,
) <-chan error {
	errCh := make(chan error, 1)

	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(SpamMessagesTotal, SpamBatchesTotal)

	go func() {
		srvErr := http.ListenAndServe(prometheusPort, nil)
		if srvErr != nil {
			logger.Error(map[string]interface{}{
				"message": srvErr.Error(),
				"error":   true,
			})
			errCh <- srvErr
		}
		close(errCh)
	}()

	return errCh
}
