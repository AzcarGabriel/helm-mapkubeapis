package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

type ManifestYaml struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App      string `yaml:"app"`
			Track    string `yaml:"track"`
			Tier     string `yaml:"tier"`
			Chart    string `yaml:"chart"`
			Release  string `yaml:"release"`
			Heritage string `yaml:"heritage"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App      string `yaml:"app"`
				Track    string `yaml:"track"`
				Tier     string `yaml:"tier"`
				Chart    string `yaml:"chart"`
				Release  string `yaml:"release"`
				Heritage string `yaml:"heritage"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App      string `yaml:"app"`
					Track    string `yaml:"track"`
					Tier     string `yaml:"tier"`
					Chart    string `yaml:"chart"`
					Release  string `yaml:"release"`
					Heritage string `yaml:"heritage"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				ImagePullSecrets []struct {
					Name string `yaml:"name"`
				} `yaml:"imagePullSecrets"`
				Containers []struct {
					Name            string `yaml:"name"`
					Image           string `yaml:"image"`
					ImagePullPolicy string `yaml:"imagePullPolicy"`
					Env             []struct {
						Name      string      `yaml:"name"`
						Value     interface{} `yaml:"value,omitempty"`
						ValueFrom struct {
							ConfigMapKeyRef struct {
								Name string `yaml:"name"`
								Key  string `yaml:"key"`
							} `yaml:"configMapKeyRef"`
						} `yaml:"valueFrom,omitempty"`
					} `yaml:"env"`
					Ports []struct {
						Name          string `yaml:"name"`
						ContainerPort int    `yaml:"containerPort"`
					} `yaml:"ports"`
					LivenessProbe struct {
						HTTPGet struct {
							Path string `yaml:"path"`
							Port int    `yaml:"port"`
						} `yaml:"httpGet"`
						InitialDelaySeconds int `yaml:"initialDelaySeconds"`
						TimeoutSeconds      int `yaml:"timeoutSeconds"`
					} `yaml:"livenessProbe"`
					ReadinessProbe struct {
						HTTPGet struct {
							Path string `yaml:"path"`
							Port int    `yaml:"port"`
						} `yaml:"httpGet"`
						InitialDelaySeconds int `yaml:"initialDelaySeconds"`
						TimeoutSeconds      int `yaml:"timeoutSeconds"`
					} `yaml:"readinessProbe"`
					VolumeMounts []struct {
						Name      string `yaml:"name"`
						MountPath string `yaml:"mountPath"`
						ReadOnly  bool   `yaml:"readOnly"`
					} `yaml:"volumeMounts"`
					Resources struct {
						Requests struct {
							CPU    string `yaml:"cpu"`
							Memory string `yaml:"memory"`
						} `yaml:"requests"`
					} `yaml:"resources"`
				} `yaml:"containers"`
				Volumes []struct {
					Name   string `yaml:"name"`
					Secret struct {
						SecretName string `yaml:"secretName"`
					} `yaml:"secret,omitempty"`
					ConfigMap struct {
						Name string `yaml:"name"`
					} `yaml:"configMap,omitempty"`
				} `yaml:"volumes"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func main() {
	varyaml := "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: testing\n  labels:\n    app: testing\n    track: \"stable\"\n    tier: \"web\"\n    chart: \"startnet-auto-deploy-0.11.0\"\n    release: testing\n    heritage: Tiller\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: testing\n      track: \"stable\"\n      tier: \"web\"\n      chart: \"startnet-auto-deploy-0.1.0\"\n      release: testing\n      heritage: Tiller\n  template:\n    metadata:\n      labels:\n        app: testing\n        track: \"stable\"\n        tier: \"web\"\n        chart: \"startnet-auto-deploy-0.1.0\"\n        release: testing\n        heritage: Tiller\n    spec:\n      imagePullSecrets:\n          - name: gitlab-registry\n\n      containers:\n      - name: startnet-auto-deploy\n        image: \"registry.gitlab.com/startnet-spa/subscription/flipcar-api/develop:16792190a806f171a7a2a71a8265100bee38bc86\"\n        imagePullPolicy: IfNotPresent\n        env:\n        - name: DATABASE_URL\n          value:\n        - name: REDISHOST\n          valueFrom:\n            configMapKeyRef:\n              name: redishost\n              key: REDISHOST\n        - name: SQLHOST\n          valueFrom:\n            configMapKeyRef:\n              name: sqlhost\n              key: SQLHOST\n        - name: CI_ENVIRONMENT\n          value: \"dev\"\n        - name: GOOGLE_APPLICATION_CREDENTIALS\n          value: /var/secrets/google/credentials.json\n        - name: ENVIRONMENT_URL\n          value: \"http://flipcar-api-testing.flipcar.cl\"\n        - name: APP_FULL_NAME\n          value: \"testing-startnet-auto-deploy\"\n        ports:\n        - name: \"web\"\n          containerPort: 5000\n        livenessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 300\n          timeoutSeconds: 15\n        readinessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 5\n          timeoutSeconds: 3\n        volumeMounts:\n        - name: parameters\n          mountPath: /tmp/parameters\n          readOnly: true\n        - name: google-application-credentials\n          mountPath: /var/secrets/google\n          readOnly: true\n        resources:\n            requests:\n              cpu: 100m\n              memory: 256Mi\n\n      volumes:\n        - name: google-application-credentials\n          secret:\n            secretName: google-application-credentials\n        - name: parameters\n          configMap:\n              name: testing-parameters"

	var yamlConfig ManifestYaml
	err := yaml.Unmarshal([]byte(varyaml), &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	// fmt.Printf("%v\n", yamlConfig.Metadata.Labels)
	// yamlConfig.Spec.Selector.MatchLabels = yamlConfig.Metadata.Labels
	// fmt.Printf("%v\n", yamlConfig.Spec.Selector.MatchLabels)
	// fmt.Printf("%v\n", yamlConfig.Spec.Template.Metadata.Labels)

	d, err := yaml.Marshal(&yamlConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}