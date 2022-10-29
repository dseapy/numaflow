package v1alpha1

type KafkaBufferService struct {
	// External holds an External Kafka config
	External *KafkaConfig `json:"external,omitempty" protobuf:"bytes,1,opt,name=external"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers,omitempty" protobuf:"bytes,1,rep,name=brokers"`
	// TLS user to configure TLS connection for kafka broker
	// TLS.enable=true default for TLS.
	// +optional
	TLS *TLS `json:"tls" protobuf:"bytes,2,opt,name=tls"`
	// +optional
	Config string `json:"config,omitempty" protobuf:"bytes,3,opt,name=config"`
}
