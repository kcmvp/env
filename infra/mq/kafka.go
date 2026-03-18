package mq

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"sync"

	"github.com/kcmvp/env/internal"
	"github.com/samber/lo"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/aws"
)

var (
	kOnce    sync.Once
	_cluster *kgo.Client
)

func Kafka() *kgo.Client {
	kOnce.Do(func() {
		var err error
		profile := internal.MstProfile()
		brokers := profile.GetStringSlice("kafka.brokers")
		topics := profile.GetStringSlice("kafka.topics")
		lo.Assert(len(topics) > 0, "kafka topics must be configured")
		group := profile.GetString("kafka.group_id")
		opts := []kgo.Opt{
			kgo.SeedBrokers(brokers...),
			kgo.ConsumeTopics(topics...),
			kgo.ConsumerGroup(group),
		}
		// support aws wsk
		accessKye := profile.GetString("kafka.accessKey")
		secretKey := profile.GetString("kafka.secretKey")
		if len(accessKye) > 0 && len(secretKey) > 0 {
			mechanism := aws.ManagedStreamingIAM(func(ctx context.Context) (aws.Auth, error) {
				return aws.Auth{
					AccessKey: accessKye,
					SecretKey: secretKey,
				}, nil
			})
			opts = append(opts, kgo.DialTLSConfig(new(tls.Config)), kgo.SASL(mechanism))
		}
		_cluster, err = kgo.NewClient(opts...)
		if err != nil {
			slog.Error("failed to create Kafka client", "error", err)
			os.Exit(1)
		}
	})
	return _cluster
}
