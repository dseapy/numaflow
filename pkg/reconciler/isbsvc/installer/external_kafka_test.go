package installer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestExternalKafkaInstallation(t *testing.T) {
	t.Run("bad installation", func(t *testing.T) {
		badIsbs := testExternalKafkaIsbSvc.DeepCopy()
		badIsbs.Spec.Kafka = nil
		installer := &externalKafkaInstaller{
			isbs:   badIsbs,
			logger: zaptest.NewLogger(t).Sugar(),
		}
		_, err := installer.Install(context.TODO())
		assert.Error(t, err)
	})

	t.Run("good installation", func(t *testing.T) {
		goodIsbs := testExternalKafkaIsbSvc.DeepCopy()
		installer := &externalKafkaInstaller{
			isbs:   goodIsbs,
			logger: zaptest.NewLogger(t).Sugar(),
		}
		c, err := installer.Install(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, c.Kafka)
		assert.Equal(t, len(goodIsbs.Spec.Kafka.External.Brokers), len(c.Kafka.Brokers))
		assert.Equal(t, goodIsbs.Spec.Kafka.External.Brokers[0], c.Kafka.Brokers[0])
		assert.Equal(t, goodIsbs.Spec.Kafka.External.Brokers[1], c.Kafka.Brokers[1])
	})
}

func TestExternalKafkaUninstallation(t *testing.T) {
	obj := testExternalKafkaIsbSvc.DeepCopy()
	installer := &externalKafkaInstaller{
		isbs:   obj,
		logger: zaptest.NewLogger(t).Sugar(),
	}
	err := installer.Uninstall(context.TODO())
	assert.NoError(t, err)
}
