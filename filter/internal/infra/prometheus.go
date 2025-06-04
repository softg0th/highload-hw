package infra

import (
	"context"
	"github.com/KaranJagtiani/go-logstash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"time"
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
	CPUUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "System-wide CPU usage percentage",
		})
	MemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Used system memory in bytes",
		})
)

func StartPrometheus(
	ctx context.Context,
	prometheusPort string,
	logger *logstash_logger.Logstash,
) <-chan error {
	errCh := make(chan error, 1)

	prometheus.MustRegister(SpamMessagesTotal, SpamBatchesTotal, CPUUsage, MemoryUsage)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				percent, err := cpu.Percent(0, false)
				if err == nil && len(percent) > 0 {
					CPUUsage.Set(percent[0])
				}

				vmem, err := mem.VirtualMemory()
				if err == nil {
					MemoryUsage.Set(float64(vmem.Used))
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(prometheusPort, nil); err != nil {
			logger.Error(map[string]interface{}{
				"message": err.Error(),
				"error":   true,
			})
			errCh <- err
		}
		close(errCh)
	}()

	return errCh
}
