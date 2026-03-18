package kafka

import (
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kfake"
)

var (
	startOnce sync.Once
	closeOnce sync.Once
	_cluster  *kfake.Cluster
	startErr  error
)

// Start starts a mock Kafka cluster. It's safe to call multiple times; subsequent
// calls will be no-ops. This Start is parameterless and will start three brokers
// listening on ports 29092,29093,29094 and enable auto topic creation.
func Start() error {
	startOnce.Do(func() {
		// Start 3 brokers on fixed ports and allow auto topic creation.
		_cluster, startErr = kfake.NewCluster(
			kfake.Ports(29092, 29093, 29094),
			kfake.AllowAutoTopicCreation(),
		)
		if startErr == nil {
			slog.Info("Mock Kafka cluster started", "brokers", _cluster.ListenAddrs())
		}
	})
	return startErr
}

// Close stops the mock Kafka cluster. It's safe to call multiple times.
// Close stops the mock Kafka cluster. It's safe to call multiple times.
// Close is idempotent (will only run once).
func Close() {
	closeOnce.Do(func() {
		if _cluster == nil {
			return
		}
		_cluster.Close()
		_cluster = nil
		slog.Info("Mock Kafka cluster stopped")
	})
}

// Brokers returns the listen addresses of the running mock cluster. If the
// cluster is not started it returns an empty slice.
func Brokers() []string {
	if _cluster == nil {
		return nil
	}
	return _cluster.ListenAddrs()
}
