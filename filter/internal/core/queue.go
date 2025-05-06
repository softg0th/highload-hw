package core

import (
	"context"
	"filter/internal/infra"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"sync"
	"time"
)

var (
	spamBuffer     = make([][]byte, 0, 100)
	spamBufferLock sync.Mutex
)

var spamProcessor func(ctx context.Context, objectName string, data [][]byte) error

func SetSpamProcessor(process func(ctx context.Context, objectName string, data [][]byte) error) {
	spamProcessor = process
}

func AddSpamMessage(msg []byte) {
	spamBufferLock.Lock()
	spamBuffer = append(spamBuffer, msg)
	spamBufferLock.Unlock()
}

func StartSpamBatchJob(interval time.Duration, logger *logstash_logger.Logstash) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			var batch [][]byte

			spamBufferLock.Lock()
			if len(spamBuffer) > 0 {
				batch = make([][]byte, len(spamBuffer))
				copy(batch, spamBuffer)
				spamBuffer = spamBuffer[:0]
			}
			spamBufferLock.Unlock()
			currentTimestamp := time.Now()
			timestampStr := currentTimestamp.Format(time.RFC3339)
			if len(batch) > 0 && spamProcessor != nil {
				logger.Info(map[string]interface{}{
					"message": "Start spam batch job",
					"error":   false,
				})

				err := spamProcessor(context.Background(), timestampStr, batch)
				if err != nil {
					logger.Error(map[string]interface{}{
						"message":           "Spam batch job error",
						"error":             true,
						"error_description": err.Error(),
					})
				}

				infra.SpamBatchesTotal.Inc()
			}
		}
	}()
}
