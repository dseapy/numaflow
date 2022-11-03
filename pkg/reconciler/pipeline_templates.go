package reconciler

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/spf13/viper"
)

// PipelineTemplates contains default values for pipeline components, if provided
// intended to be populated from the configmap attached to the controller manager.
type PipelineTemplates struct {
	DaemonTemplate *dfv1.DaemonTemplate `json:"daemon"`
	JobTemplate    *dfv1.JobTemplate    `json:"job"`
	VertexTemplate *VertexTemplate      `json:"vertex"`
}

type VertexTemplate struct {
	PodTemplate           *dfv1.AbstractPodTemplate `json:"podTemplate"`
	ContainerTemplate     *dfv1.ContainerTemplate   `json:"containerTemplate"`
	InitContainerTemplate *dfv1.ContainerTemplate   `json:"initContainerTemplate"`
}

func LoadPipelineTemplates(onErrorReloading func(error)) (*PipelineTemplates, error) {
	v := viper.New()
	v.SetConfigName("pipeline-templates")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/numaflow")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to load pipeline-templates file. %w", err)
		}
	}
	r := &PipelineTemplates{}
	err = v.Unmarshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal pipeline-templates file. %w", err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//TODO: is this thread safe?  also in config.go. if not, add lock or just require restart?
		//  If removing, don't need to pass entire PipelineTemplates to vertex controller, just VertexTemplate needed
		err = v.Unmarshal(r)
		if err != nil {
			onErrorReloading(err)
		}
	})
	return r, nil
}
