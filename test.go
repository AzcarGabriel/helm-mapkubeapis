package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
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
	varyaml := "---\n# Source: startnet-auto-deploy/templates/service.yaml\napiVersion: v1\nkind: Service\nmetadata:\n  name: testing-startnet-auto-deploy\n  labels:\n    app: testing\n    chart: \"startnet-auto-deploy-0.1.0\"\n    release: testing\n    heritage: Tiller\nspec:\n  type: ClusterIP\n  ports:\n  - port: 5000\n    targetPort: 5000\n    protocol: TCP\n    name: web\n  selector:\n    app: testing\n    release: testing\n    tier: \"web\"\n---\n# Source: startnet-auto-deploy/templates/deployment.yaml\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: testing\n  labels:\n    app: testing\n    track: \"stable\"\n    tier: \"web\"\n    chart: \"startnet-auto-deploy-0.1.0\"\n    release: testing\n    heritage: Tiller\nspec:\n  replicas: 1\n  template:\n    metadata:\n      labels:\n        app: testing\n        track: \"stable\"\n        tier: \"web\"\n        release: testing\n    spec:\n      imagePullSecrets:\n          - name: gitlab-registry\n      containers:\n      - name: startnet-auto-deploy\n        image: \"registry.gitlab.com/startnet-spa/subscription/flipcar-api/develop:16792190a806f171a7a2a71a8265100bee38bc86\"\n        imagePullPolicy: IfNotPresent\n        env:\n        - name: DATABASE_URL\n          value:\n        - name: REDISHOST\n          valueFrom:\n            configMapKeyRef:\n              name: redishost\n              key: REDISHOST\n        - name: SQLHOST\n          valueFrom:\n            configMapKeyRef:\n              name: sqlhost\n              key: SQLHOST\n        - name: CI_ENVIRONMENT\n          value: \"dev\"\n        - name: GOOGLE_APPLICATION_CREDENTIALS\n          value: /var/secrets/google/credentials.json\n        - name: ENVIRONMENT_URL\n          value: \"http://flipcar-api-testing.flipcar.cl\"\n        - name: APP_FULL_NAME\n          value: \"testing-startnet-auto-deploy\"\n        ports:\n        - name: \"web\"\n          containerPort: 5000\n        livenessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 300\n          timeoutSeconds: 15\n        readinessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 5\n          timeoutSeconds: 3\n        volumeMounts:\n        - name: parameters\n          mountPath: /tmp/parameters\n          readOnly: true\n        - name: google-application-credentials\n          mountPath: /var/secrets/google\n          readOnly: true\n        resources:\n            requests:\n              cpu: 100m\n              memory: 256Mi\n      volumes:\n        - name: google-application-credentials\n          secret:\n            secretName: google-application-credentials\n        - name: parameters\n          configMap:\n              # Provide the name of the ConfigMap containing the files you want\n              # to add to the container\n              name: testing-parameters\n---\n# Source: startnet-auto-deploy/templates/jobs.yaml\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: \"testing-queue-consumer\"\n  labels:\n    chart: \"startnet-auto-deploy-0.1.0\"\nspec:\n  replicas: 1\n  template:\n    metadata:\n      labels:\n        app: testing\n    spec:\n      imagePullSecrets:\n          - name: gitlab-registry\n      containers:\n      - name: queue-consumer\n        image: \"registry.gitlab.com/startnet-spa/subscription/flipcar-api/develop:16792190a806f171a7a2a71a8265100bee38bc86\"\n        imagePullPolicy: IfNotPresent\n        env:\n          - name: DATABASE_URL\n            value:\n          - name: REDISHOST\n            valueFrom:\n              configMapKeyRef:\n                name: redishost\n                key: REDISHOST\n          - name: SQLHOST\n            valueFrom:\n              configMapKeyRef:\n                name: sqlhost\n                key: SQLHOST\n          - name: CI_ENVIRONMENT\n            value: \"dev\"\n          - name: GOOGLE_APPLICATION_CREDENTIALS\n            value: /var/secrets/google/credentials.json\n        args:\n          - consume\n          - -vvv\n          - enqueue\n        volumeMounts:\n        - name: parameters\n          mountPath: /tmp/parameters\n          readOnly: true\n        - name: google-application-credentials\n          mountPath: /var/secrets/google\n          readOnly: true\n        resources:\n          requests:\n            cpu: 100m\n            memory: 256Mi\n      volumes:\n        - name: google-application-credentials\n          secret:\n            secretName: google-application-credentials\n        - name: parameters\n          configMap:\n            # Provide the name of the ConfigMap containing the files you want\n            # to add to the container\n            name: \"testing-parameters\"\n---\n# Source: startnet-auto-deploy/templates/ingress.yaml\napiVersion: networking.k8s.io/v1beta1\nkind: Ingress\nmetadata:\n  name: testing-startnet-auto-deploy\n  labels:\n    app: testing\n    chart: \"startnet-auto-deploy-0.1.0\"\n    release: testing\n    heritage: Tiller\n  annotations:\n    kubernetes.io/tls-acme: \"true\"\n    kubernetes.io/ingress.class: \"nginx\"\n    nginx.ingress.kubernetes.io/proxy-body-size: 20m\nspec:\n  tls:\n  - hosts:\n    - \"flipcar-api-testing.flipcar.cl\"\n    secretName: testing-startnet-auto-deploy-tls\n  rules:\n  - host: \"flipcar-api-testing.flipcar.cl\"\n    http:\n      paths:\n      - path: /\n        backend:\n          serviceName: testing-startnet-auto-deploy\n          servicePort: 5000"

	ans := ""
	parts := strings.Split(varyaml, "---")
	for _, s := range parts {
		var yamlConfig ManifestYaml
		err := yaml.Unmarshal([]byte(s), &yamlConfig)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}

		if yamlConfig.Kind == "Deployment" && yamlConfig.Metadata.Name == "testing" {
			var manifestYaml ManifestYaml
			err = yaml.Unmarshal([]byte(s), &manifestYaml)
			if err != nil {
				log.Printf("Error parsing YAML file: %s\n", err)
			}
			manifestYaml.Spec.Selector.MatchLabels = manifestYaml.Metadata.Labels
			manifestYaml.Spec.Template.Metadata.Labels = manifestYaml.Metadata.Labels

			yamlString, err := yaml.Marshal(&manifestYaml)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			ans += string(yamlString)
		} else {
			ans += s
		}
	}

	fmt.Printf("%s", ans)
}