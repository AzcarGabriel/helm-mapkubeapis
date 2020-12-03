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

type Labels struct {
	App      string `yaml:"app"`
	Track    string `yaml:"track"`
	Tier     string `yaml:"tier"`
	Chart    string `yaml:"chart"`
	Release  string `yaml:"release"`
	Heritage string `yaml:"heritage"`
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
	var labels = constructLabels(modifiedManifest)
	for _, s := range parts {
		var yamlConfig ManifestYaml
		err := yaml.Unmarshal([]byte(s), &yamlConfig)
		if err != nil {
			log.Printf("Error parsing YAML file: %s\n", err)
		}

		if yamlConfig.Kind == "Deployment" {
			var manifestYaml ManifestYaml
			err = yaml.Unmarshal([]byte(s), &manifestYaml)
			if err != nil {
				log.Printf("Error parsing YAML file: %s\n", err)
			}
			manifestYaml.Metadata.Labels = labels
			manifestYaml.Spec.Selector.MatchLabels = labels
			manifestYaml.Spec.Template.Metadata.Labels = labels

			yamlString, err := yaml.Marshal(&manifestYaml)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			finalManifest += "---\n" + string(yamlString)
		} else {
			finalManifest += "---" + s
		}
	}
	log.Printf(finalManifest)
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

// constructLabels returns labels variables for Deployment kind
func constructLabels(manifest string) Labels {
	// Variables
	parts := strings.Split(manifest, "---")
	var labels Labels
	labels.Tier = "web"
	labels.Track = "stable"
	labels.Heritage = "Tiller"

	// Revision
	for _, s := range parts {
		var yamlConfig ManifestYaml
		err := yaml.Unmarshal([]byte(s), &yamlConfig)
		if err != nil {
			log.Printf("Error parsing YAML file: %s\n", err)
		}

		if yamlConfig.Kind == "Deployment" {
			var manifestYaml ManifestYaml
			err = yaml.Unmarshal([]byte(s), &manifestYaml)
			if err != nil {
				log.Printf("Error parsing YAML file: %s\n", err)
			}

			// Relleno de labels
			if labels.App == "" {
				labels.App = manifestYaml.Metadata.Labels.App
			}
			if labels.Chart == "" {
				labels.Chart = manifestYaml.Metadata.Labels.Chart
			}
			if labels.Release == "" {
				labels.Release = manifestYaml.Metadata.Labels.Release
			}
		}
	}

	return labels
}
