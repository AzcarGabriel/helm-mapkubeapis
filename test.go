package main

import (
	"fmt"
	"github.com/hickeyma/helm-mapkubeapis/pkg/mapping"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strings"
)

// Yaml del manifiesto
type DeploymentYaml struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels Labels `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels
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

type Labels struct {
	App      string `yaml:"app,omitempty"`
	Track    string `yaml:"track,omitempty"`
	Tier     string `yaml:"tier,omitempty"`
	Chart    string `yaml:"chart,omitempty"`
	Release  string `yaml:"release,omitempty"`
}

type IngressYaml struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App      string `yaml:"app"`
			Chart    string `yaml:"chart"`
			Release  string `yaml:"release"`
			Heritage string `yaml:"heritage"`
			AppKubernetesIoManagedBy string `yaml:"app.kubernetes.io/managed-by,omitempty"`
		} `yaml:"labels"`
		Annotations `yaml:"annotations,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		TLS []struct {
			Hosts      []string `yaml:"hosts"`
			SecretName string   `yaml:"secretName"`
		} `yaml:"tls"`
		Rules []struct {
			Host string `yaml:"host"`
			HTTP struct {
				Paths []struct {
					Path    string `yaml:"path"`
					Backend struct {
						ServiceName string `yaml:"serviceName"`
						ServicePort int    `yaml:"servicePort"`
					} `yaml:"backend"`
				} `yaml:"paths"`
			} `yaml:"http"`
		} `yaml:"rules"`
	} `yaml:"spec"`
}

type Annotations struct {
	KubernetesIoTLSAcme                   string `yaml:"kubernetes.io/tls-acme,omitempty"`
	KubernetesIoIngressClass              string `yaml:"kubernetes.io/ingress.class,omitempty"`
	NginxIngressKubernetesIoProxyBodySize string `yaml:"nginx.ingress.kubernetes.io/proxy-body-size,omitempty"`
	MetaHelmShReleaseName                 string `yaml:"meta.helm.sh/release-name,omitempty"`
}

func main() {
	origManifest := "---\napiVersion: v1\nkind: Service\nmetadata:\n  name: testing-startnet-auto-deploy\n  labels:\n    app: testing\n    chart: \"startnet-auto-deploy-0.1.0\"\n    release: testing\n    heritage: Tiller\nspec:\n  type: ClusterIP\n  ports:\n  - port: 5000\n    targetPort: 5000\n    protocol: TCP\n    name: web\n  selector:\n    app: testing\n    release: testing\n    tier: \"web\"\n---\napiVersion: extensions/v1beta1\nkind: Deployment\nmetadata:\n  name: testing\n  labels:\n    app: testing\n    track: stable\n    tier: web\n    chart: startnet-auto-deploy-0.1.0\n    release: testing\n    heritage: Tiller\nspec:\n  replicas: 1\n  selector:\n    matchlabels:\n      app: testing\n      track: stable\n      tier: web\n      release: testing\n  template:\n    metadata:\n      labels:\n        app: testing\n        track: stable\n        tier: web\n        release: testing\n    spec:\n      imagePullSecrets:\n      - name: gitlab-registry\n      containers:\n      - name: startnet-auto-deploy\n        image: registry.gitlab.com/startnet-spa/autosimple/autosimple-api/develop:fc721874f33ec367549c6dc34e9a49446d31c5ed\n        imagePullPolicy: IfNotPresent\n        env:\n        - name: DATABASE_URL\n        - name: REDISHOST\n          valueFrom:\n            configMapKeyRef:\n              name: redishost\n              key: REDISHOST\n        - name: SQLHOST\n          valueFrom:\n            configMapKeyRef:\n              name: sqlhost\n              key: SQLHOST\n        - name: CI_ENVIRONMENT\n          value: dev\n        - name: GOOGLE_APPLICATION_CREDENTIALS\n          value: /var/secrets/google/credentials.json\n        - name: ENVIRONMENT_URL\n          value: http://autosimple-api-testing.credisimple.cl\n        ports:\n        - name: web\n          containerPort: 5000\n        livenessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 300\n          timeoutSeconds: 15\n        readinessProbe:\n          httpGet:\n            path: /health/\n            port: 5000\n          initialDelaySeconds: 5\n          timeoutSeconds: 3\n        volumeMounts:\n        - name: parameters\n          mountPath: /tmp/parameters\n          readOnly: true\n        - name: google-application-credentials\n          mountPath: /var/secrets/google\n          readOnly: true\n        resources:\n          requests:\n            cpu: 100m\n            memory: 256Mi\n      volumes:\n      - name: google-application-credentials\n        secret:\n          secretName: google-application-credentials\n      - name: parameters\n        configMap:\n          name: testing-parameters\n---\napiVersion: extensions/v1beta1\nkind: Deployment\nmetadata:\n  name: testing-queue-consumer\n  labels:\n    chart: startnet-auto-deploy-0.1.0\nspec:\n  replicas: 1\n  selector:\n    matchlabels:\n      app: testing\n  template:\n    metadata:\n      labels:\n        app: testing\n    spec:\n      imagePullSecrets:\n      - name: gitlab-registry\n      containers:\n      - name: queue-consumer\n        image: registry.gitlab.com/startnet-spa/autosimple/autosimple-api/develop:fc721874f33ec367549c6dc34e9a49446d31c5ed\n        imagePullPolicy: IfNotPresent\n        env:\n        - name: DATABASE_URL\n        - name: REDISHOST\n          valueFrom:\n            configMapKeyRef:\n              name: redishost\n              key: REDISHOST\n        - name: SQLHOST\n          valueFrom:\n            configMapKeyRef:\n              name: sqlhost\n              key: SQLHOST\n        - name: CI_ENVIRONMENT\n          value: dev\n        - name: GOOGLE_APPLICATION_CREDENTIALS\n          value: /var/secrets/google/credentials.json\n        ports: []\n        livenessProbe:\n          httpGet:\n            path: \"\"\n            port: 0\n          initialDelaySeconds: 0\n          timeoutSeconds: 0\n        readinessProbe:\n          httpGet:\n            path: \"\"\n            port: 0\n          initialDelaySeconds: 0\n          timeoutSeconds: 0\n        volumeMounts:\n        - name: parameters\n          mountPath: /tmp/parameters\n          readOnly: true\n        - name: google-application-credentials\n          mountPath: /var/secrets/google\n          readOnly: true\n        resources:\n          requests:\n            cpu: 100m\n            memory: 256Mi\n      volumes:\n      - name: google-application-credentials\n        secret:\n          secretName: google-application-credentials\n      - name: parameters\n        configMap:\n          name: testing-parameters\n---\n# Source: startnet-auto-deploy/templates/cronjob.yaml\napiVersion: batch/v1beta1\nkind: CronJob\nmetadata:\n  name: \"testing-get-uf\"\n  labels:\n    chart: \"startnet-auto-deploy-0.1.0\"\nspec:\n  schedule: \"1 0 * * *\"\n  concurrencyPolicy: Forbid\n  successfulJobsHistoryLimit: 1\n  failedJobsHistoryLimit: 1\n  jobTemplate:\n    spec:\n      template:\n        metadata:\n          labels:\n            app: testing\n            cron: get-uf\n        spec:\n          imagePullSecrets:\n            - name: gitlab-registry\n          containers:\n          - name: get-uf\n            image: \"registry.gitlab.com/startnet-spa/autosimple/autosimple-api/develop:fc721874f33ec367549c6dc34e9a49446d31c5ed\"\n            imagePullPolicy: IfNotPresent\n            env:\n            - name: DATABASE_URL\n              value:\n            - name: REDISHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: redishost\n                      key: REDISHOST\n            - name: SQLHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: sqlhost\n                      key: SQLHOST\n            - name: CI_ENVIRONMENT\n              value: \"dev\"\n            - name: GOOGLE_APPLICATION_CREDENTIALS\n              value: /var/secrets/google/credentials.json\n            args:\n              - raw_console\n              - -vvv\n              - autosimple:get-uf\n            volumeMounts:\n            - name: parameters\n              mountPath: /tmp/parameters\n              readOnly: true\n            - name: google-application-credentials\n              mountPath: /var/secrets/google\n              readOnly: true\n            resources:\n              requests:\n                cpu: 100m\n                memory: 256Mi\n          volumes:\n          - name: google-application-credentials\n            secret:\n              secretName: google-application-credentials\n          - name: parameters\n            configMap:\n              # Provide the name of the ConfigMap containing the files you want\n              # to add to the container\n              name: \"testing-parameters\"\n          restartPolicy: Never\n---\n# Source: startnet-auto-deploy/templates/cronjob.yaml\napiVersion: batch/v1beta1\nkind: CronJob\nmetadata:\n  name: \"testing-campaign-decoder\"\n  labels:\n    chart: \"startnet-auto-deploy-0.1.0\"\nspec:\n  schedule: \"*/5 * * * *\"\n  concurrencyPolicy: Forbid\n  successfulJobsHistoryLimit: 1\n  failedJobsHistoryLimit: 1\n  jobTemplate:\n    spec:\n      template:\n        metadata:\n          labels:\n            app: testing\n            cron: campaign-decoder\n        spec:\n          imagePullSecrets:\n            - name: gitlab-registry\n          containers:\n          - name: campaign-decoder\n            image: \"registry.gitlab.com/startnet-spa/autosimple/autosimple-api/develop:fc721874f33ec367549c6dc34e9a49446d31c5ed\"\n            imagePullPolicy: IfNotPresent\n            env:\n            - name: DATABASE_URL\n              value:\n            - name: REDISHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: redishost\n                      key: REDISHOST\n            - name: SQLHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: sqlhost\n                      key: SQLHOST\n            - name: CI_ENVIRONMENT\n              value: \"dev\"\n            - name: GOOGLE_APPLICATION_CREDENTIALS\n              value: /var/secrets/google/credentials.json\n            args:\n              - raw_console\n              - -vvv\n              - autosimple:campaign-decoder\n            volumeMounts:\n            - name: parameters\n              mountPath: /tmp/parameters\n              readOnly: true\n            - name: google-application-credentials\n              mountPath: /var/secrets/google\n              readOnly: true\n            resources:\n              requests:\n                cpu: 100m\n                memory: 256Mi\n          volumes:\n          - name: google-application-credentials\n            secret:\n              secretName: google-application-credentials\n          - name: parameters\n            configMap:\n              # Provide the name of the ConfigMap containing the files you want\n              # to add to the container\n              name: \"testing-parameters\"\n          restartPolicy: Never\n---\n# Source: startnet-auto-deploy/templates/cronjob.yaml\napiVersion: batch/v1beta1\nkind: CronJob\nmetadata:\n  name: \"testing-wake-up-offers\"\n  labels:\n    chart: \"startnet-auto-deploy-0.1.0\"\nspec:\n  schedule: \"*/1 * * * *\"\n  concurrencyPolicy: Forbid\n  successfulJobsHistoryLimit: 1\n  failedJobsHistoryLimit: 1\n  jobTemplate:\n    spec:\n      template:\n        metadata:\n          labels:\n            app: testing\n            cron: wake-up-offers\n        spec:\n          imagePullSecrets:\n            - name: gitlab-registry\n          containers:\n          - name: wake-up-offers\n            image: \"registry.gitlab.com/startnet-spa/autosimple/autosimple-api/develop:fc721874f33ec367549c6dc34e9a49446d31c5ed\"\n            imagePullPolicy: IfNotPresent\n            env:\n            - name: DATABASE_URL\n              value:\n            - name: REDISHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: redishost\n                      key: REDISHOST\n            - name: SQLHOST\n              valueFrom:\n                  configMapKeyRef:\n                      name: sqlhost\n                      key: SQLHOST\n            - name: CI_ENVIRONMENT\n              value: \"dev\"\n            - name: GOOGLE_APPLICATION_CREDENTIALS\n              value: /var/secrets/google/credentials.json\n            args:\n              - raw_console\n              - -vvv\n              - autosimple:wake-up-offers\n            volumeMounts:\n            - name: parameters\n              mountPath: /tmp/parameters\n              readOnly: true\n            - name: google-application-credentials\n              mountPath: /var/secrets/google\n              readOnly: true\n            resources:\n              requests:\n                cpu: 100m\n                memory: 256Mi\n          volumes:\n          - name: google-application-credentials\n            secret:\n              secretName: google-application-credentials\n          - name: parameters\n            configMap:\n              # Provide the name of the ConfigMap containing the files you want\n              # to add to the container\n              name: \"testing-parameters\"\n          restartPolicy: Never\n---\n# Source: startnet-auto-deploy/templates/ingress.yaml\napiVersion: extensions/v1beta1\nkind: Ingress\nmetadata:\n  name: testing-startnet-auto-deploy\n  labels:\n    app: testing\n    chart: \"startnet-auto-deploy-0.1.0\"\n    release: testing\n    heritage: Tiller\n  annotations:\n    kubernetes.io/tls-acme: \"true\"\n    kubernetes.io/ingress.class: \"nginx\"\n    nginx.ingress.kubernetes.io/proxy-body-size: 20m\nspec:\n  tls:\n  - hosts:\n    - \"autosimple-api-testing.credisimple.cl\"\n    secretName: testing-startnet-auto-deploy-tls\n  - hosts:\n    - \"autosimple-api-testing.autosimple.cl\"\n    secretName: autosimple-api-testing-autosimple-cl-tls\n  rules:\n  - host: \"autosimple-api-testing.autosimple.cl\"\n    http:\n      paths:\n      - path: /\n        backend:\n          serviceName: testing-startnet-auto-deploy\n          servicePort: 5000\n  - host: \"autosimple-api-testing.credisimple.cl\"\n    http:\n      paths:\n      - path: /\n        backend:\n          serviceName: testing-startnet-auto-deploy\n          servicePort: 5000\n"
	var modifiedManifest = origManifest
	var err error
	var mapMetadata *mapping.Metadata

	// Load the mapping data
	if mapMetadata, err = mapping.LoadMapfile("config/Map.yaml"); err != nil {
		return
	}

	// Check for deprecated or removed APIs and map accordingly to supported versions
	for _, mapping := range mapMetadata.Mappings {
		deprecatedAPI := mapping.DeprecatedAPI
		supportedAPI := mapping.NewAPI
		var apiVersionStr string
		if mapping.DeprecatedInVersion != "" {
			apiVersionStr = mapping.DeprecatedInVersion
		} else {
			apiVersionStr = mapping.RemovedInVersion
		}

		var modManifestForAPI string
		var modified = false

		// Replace using regex
		var re = regexp.MustCompile(deprecatedAPI)
		modManifestForAPI = re.ReplaceAllString(modifiedManifest, supportedAPI)

		if modManifestForAPI != modifiedManifest {
			modified = true
			log.Printf("Found deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", deprecatedAPI, supportedAPI)
		}
		if modified {
			if semver.Compare(apiVersionStr, "v1.16") > 0 {
				log.Printf("The following API does not require mapping as the "+
					"API is not deprecated or removed in Kubernetes '%s':\n\"%s\"\n", apiVersionStr,
					deprecatedAPI)
			} else {
				modifiedManifest = modManifestForAPI
			}
		}
	}

	// Add labels variables
	finalManifest := ""
	parts := strings.Split(modifiedManifest, "---")

	// Revision
	for _, s := range parts {
		if !strings.Contains(s, "apiVersion") {
			continue
		}
		var deploymentYaml DeploymentYaml
		err := yaml.Unmarshal([]byte(s), &deploymentYaml)
		if err != nil {
			log.Printf("Error parsing YAML file: %s\n", err)
		}

		if deploymentYaml.Kind == "Deployment" {
			deploymentYaml.Spec.Selector.MatchLabels = deploymentYaml.Spec.Template.Metadata.Labels

			yamlString, err := yaml.Marshal(&deploymentYaml)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			finalManifest += "---\n" + string(yamlString)
		} else if deploymentYaml.Kind == "Ingress" {
			var ingressYaml IngressYaml
			err := yaml.Unmarshal([]byte(s), &ingressYaml)
			if err != nil {
				log.Printf("Error parsing YAML file: %s\n", err)
			}

			ingressYaml.Metadata.Labels.AppKubernetesIoManagedBy = "" // "Helm"
			ingressYaml.Metadata.Annotations.MetaHelmShReleaseName = "" // ingressYaml.Metadata.Labels.Release

			yamlString, err := yaml.Marshal(&ingressYaml)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			finalManifest += "---\n" + string(yamlString)
		} else {
			finalManifest += "---" + s
		}
	}
	fmt.Printf("%s\n", finalManifest)
}