package kafka

import (
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kfake"
)

var once sync.Once
var _cluster *kfake.Cluster

func Start() {
	once.Do(func() {
		var err error
		opts := []kfake.Opt{kfake.NumBrokers(2),
			kfake.Ports(57351, 57352)}
		_cluster, err = kfake.NewCluster(opts...)
		if err != nil {
			panic(err)
		}
		slog.Info("Mock Kafka cluster started", "brokers", _cluster.ListenAddrs())
	})
}
