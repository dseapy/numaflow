# Controller ConfigMap

The controller ConfigMap is used for controller-wide settings.

For a detailed example, please see [`numaflow-controller-config.yaml`](./numaflow-controller-config.yaml).

## Configuration Structure

The configuration should be under `controller-config.yaml` key in the ConfigMap, as a string in `yaml` format.  Pipeline templates, if provided, should be under `pipeline-templates.yaml`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: numaflow-controller-config
data:
  controller-config.yaml: |
    isbsvc:
      jetstream:
        ...
  # optional
  pipeline-templates.yaml: |
    vertex:
      ...
```

### ISB Service Configuration

One of the important configuration items in the ConfigMap is about [ISB Service](./inter-step-buffer-service.md). We currently use 3rd party technologies such as `JetStream` to implement ISB Services, if those applications have new releases, to make them available in Numaflow, the new versions need to be added in the ConfigMap.

For example, there's a new `Nats JetStream` version `x.y.x` available, a new version configuration like below needs to be added before it can be referenced in the `InterStepBufferService` spec.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: numaflow-controller-config
data:
  controller-config.yaml: |
    isbsvc:
      jetstream:
        versions:
          - version: x.y.x   # Name it whatever you want, it will be referenced in the InterStepBufferService spec.
            natsImage: nats:x.y.x
            metricsExporterImage: natsio/prometheus-nats-exporter:0.9.1
            configReloaderImage: natsio/nats-server-config-reloader:0.7.0
            startCommand: /nats-server
```

### Pipeline Templates

Pipeline Templates provide default Pipeline component configuration, reducing the need to duplicate common configuration in each Pipeline managed by the controller.

The numaflow-controller logs `Successfully loaded provided pipeline-templates file` if it detects pipeline templates.

* `daemon` has the same structure as `.spec.templates.daemon` in a Pipeline, see [Daemon customization](./pipeline-customization.md#daemon-deployment) example.
* `job` has the same structure as `.spec.templates.job` in a Pipeline, see [Job customization](./pipeline-customization.md#job) example.
* `vertex` currently has 3 possible fields:
    * `podTemplate` has the same structure as `.spec.vertices[*]` in a Pipeline, however only the fields corresponding to a Pod metadata or spec.
  Specifically, it supports the same pod metadata and spec fields that daemon and job support, see [Pipeline customization](./pipeline-customization.md) examples.
    * `containerTemplate` has the same structure as `.spec.vertices[*].containerTemplate` in a Pipeline.
    * `initContainerTemplate` has the same structure as `.spec.vertices[*].initContainerTemplate` in a Pipeline.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: numaflow-controller-config
data:
  pipeline-templates.yaml: |
    daemon:
      replicas: 2
    job:
      ttlSecondsAfterFinished: 600
    vertex:
      podTemplate:
        metadata:
          annotations:
            key1: value1
        priorityClassName: my-priority-class-name
      containerTemplate:
        resources:
          requests:
            cpu: 200m
      initContainerTemplate:
        resources:
          limits:
            memory: 256Mi
```
