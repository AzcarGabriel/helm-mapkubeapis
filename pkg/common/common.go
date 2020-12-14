/*
Copyright

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strings"

	utils "github.com/maorfr/helm-plugin-utils/pkg"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"

	"github.com/hickeyma/helm-mapkubeapis/pkg/mapping"
)

// KubeConfig are the Kubernetes configuration settings
type KubeConfig struct {
	Context string
	File    string
}

// MapOptions are the options for mapping deprecated APIs in a release
type MapOptions struct {
	DryRun           bool
	KubeConfig       KubeConfig
	MapFile          string
	ReleaseName      string
	ReleaseNamespace string
	StorageType      string
	TillerOutCluster bool
}

// Yaml del manifiesto
type ManifestYaml struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels Labels
		Annotations `yaml:"annotations,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels Labels
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels Labels
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
	AppKubernetesIoManagedBy              string `yaml:"app.kubernetes.io/managed-by,omitempty"`
	MetaHelmShReleaseName                 string `yaml:"meta.helm.sh/release-name,omitempty"`
}

// UpgradeDescription is description of why release was upgraded
const UpgradeDescription = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"

// ReplaceManifestUnSupportedAPIs returns a release manifest with deprecated or removed
// Kubernetes APIs updated to supported APIs
func ReplaceManifestUnSupportedAPIs(origManifest, mapFile string, kubeConfig KubeConfig) (string, error) {
	var modifiedManifest = origManifest
	var err error
	var mapMetadata *mapping.Metadata

	// Load the mapping data
	if mapMetadata, err = mapping.LoadMapfile(mapFile); err != nil {
		return "", errors.Wrapf(err, "Failed to load mapping file: %s", mapFile)
	}

	// get the Kubernetes server version
	kubeVersionStr, err := getKubernetesServerVersion(kubeConfig)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(kubeVersionStr) {
		return "", errors.Errorf("Failed to get Kubernetes server version")
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
		if !semver.IsValid(apiVersionStr) {
			return "", errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s", strings.ReplaceAll(deprecatedAPI, "\n", " "))
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
			if semver.Compare(apiVersionStr, kubeVersionStr) > 0 {
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

			ingressYaml.Metadata.Annotations.AppKubernetesIoManagedBy = "Helm"
			ingressYaml.Metadata.Annotations.MetaHelmShReleaseName = ingressYaml.Metadata.Labels.Release

			yamlString, err := yaml.Marshal(&ingressYaml)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			finalManifest += "---\n" + string(yamlString)
		} else {
			finalManifest += "---" + s
		}
	}
	log.Printf("%s\n", finalManifest)
	return finalManifest, nil
}

func getKubernetesServerVersion(kubeConfig KubeConfig) (string, error) {
	clientSet := utils.GetClientSetWithKubeConfig(kubeConfig.File, kubeConfig.Context)
	if clientSet == nil {
		return "", errors.Errorf("kubernetes cluster unreachable")
	}
	kubeVersion, err := clientSet.ServerVersion()
	if err != nil {
		return "", errors.Wrap(err, "kubernetes cluster unreachable")
	}
	return kubeVersion.GitVersion, nil
}