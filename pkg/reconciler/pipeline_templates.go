package reconciler

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/spf13/viper"
)

func LoadPipelineTemplates(onErrorReloading func(error)) (*dfv1.Templates, error) {
	v := viper.New()
	v.SetConfigName("pipeline-templates")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/numaflow")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to load pipeline templates file. %w", err)
		}
	}
	r := &dfv1.Templates{}
	err = v.Unmarshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal pipeline templates file. %w", err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//TODO: is this thread safe?  also in config.go. if not
		//  * is it likely to matter (config reload at same time as being used in reconcile)?
		//  * does it make sense to add lock or just require restart?
		//  If removing, don't need to pass entire Templates to vertex controller, just VertexTemplate needed
		err = v.Unmarshal(r)
		if err != nil {
			onErrorReloading(err)
		}
	})
	return r, nil
}
