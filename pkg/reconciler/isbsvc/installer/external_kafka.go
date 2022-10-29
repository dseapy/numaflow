package installer

import (
	"context"
	"fmt"

	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"go.uber.org/zap"
)

type externalKafkaInstaller struct {
	isbs   *dfv1.InterStepBufferService
	logger *zap.SugaredLogger
}

func NewExternalKafkaInstaller(isbs *dfv1.InterStepBufferService, logger *zap.SugaredLogger) Installer {
	return &externalKafkaInstaller{
		isbs:   isbs,
		logger: logger.With("isbs", isbs.Name),
	}
}

func (eki *externalKafkaInstaller) Install(ctx context.Context) (*dfv1.BufferServiceConfig, error) {
	if eki.isbs.Spec.Kafka == nil || eki.isbs.Spec.Kafka.External == nil {
		return nil, fmt.Errorf("invalid InterStepBufferService spec, no external config")
	}
	eki.isbs.Status.SetType(dfv1.ISBSvcTypeKafka)
	eki.isbs.Status.MarkConfigured()
	eki.isbs.Status.MarkDeployed()
	eki.logger.Info("Using external kafka config")
	return &dfv1.BufferServiceConfig{Kafka: eki.isbs.Spec.Kafka.External}, nil
}

func (eri *externalKafkaInstaller) Uninstall(ctx context.Context) error {
	eri.logger.Info("Nothing to uninstall")
	return nil
}
