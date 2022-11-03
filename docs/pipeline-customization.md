# Pipeline Customization

There is an optional `.spec.templates` field in the `Pipeline` resource which may be used to customize kubernetes resources owned by the Pipeline.
These override anything specified in the [Controller ConfigMap Pipeline Templates](./controller-configmap.md#pipeline-templates).

Individual vertex customization is described separately in more detail (i.e. [Environment Variables](./environment-variables.md), [Container Resources](./container-resources.md), etc.)
and take precedence over any pipeline templates.

## Component customization

The following example shows all currently supported fields. The `.spec.templates.<component>` field and all fields directly under it are optional.

```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: my-pipeline
spec:
  templates:
    # can be "daemon", "job" or "vertex"
    daemon:
      # Pod metadata
      metadata:
        labels:
          my-label-key: my-label-value
        annotations:
          my-annotation-key: my-annotation-value
      # Pod spec
      nodeSelector:
        my-node-label-key: my-node-label-value
      tolerations:
        - key: "my-example-key"
          operator: "Exists"
          effect: "NoSchedule"
      securityContext: {}
      imagePullSecrets:
        - name: regcred
      priorityClassName: my-priority-class-name
      priority: 50
      serviceAccountName: my-service-account
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app.kubernetes.io/component
                    operator: In
                    values:
                      - daemon
                  - key: numaflow.numaproj.io/pipeline-name
                    operator: In
                    values:
                      - my-pipeline
              topologyKey: kubernetes.io/hostname
      # Containers
      containerTemplate:
        env:
          - name: MY_ENV_NAME
            value: my-env-value
        resources:
          limits:
            memory: 500Mi
      initContainerTemplate:
        env:
          - name: MY_ENV_NAME
            value: my-env-value
        resources:
          limits:
            memory: 500Mi
```

## Daemon customization

In addition to the `Component customization` described above, the Pipeline daemon has the following additional fields available.

```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: my-pipeline
spec:
  templates:
    daemon:
      replicas: 3
```

## Job customization

In addition to the `Component customization` described above, Pipeline jobs have the following additional fields available.

```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: my-pipeline
spec:
  templates:
    job:
      ttlSecondsAfterFinished: 600 # numaflow defaults to 30
      backoffLimit: 5 # numaflow defaults to 20
```

## Vertex customization
* `vertex` currently has 3 possible fields:
* `podTemplate` has the same structure as `.spec.vertices[*]` in a Pipeline, however only the fields corresponding to a Pod metadata or spec.
  Specifically, it supports the same pod metadata and spec fields that daemon and job support, see [Pipeline customization](./pipeline-customization.md) examples.
    * `containerTemplate` has the same structure as `.spec.vertices[*].containerTemplate` in a Pipeline.
    * `initContainerTemplate` has the same structure as `.spec.vertices[*].initContainerTemplate` in a Pipeline.

See `Component customization` described above.