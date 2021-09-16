/*
Copyright 2018-2021 The Kubernetes Authors.

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

package system

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"k8s.io/klog/v2"

	"openshift/node-feature-discovery/source"
)

var osReleaseFields = [...]string{
	"ID",
	"VERSION_ID",
	"RHEL_VERSION",
	"OPENSHIFT_VERSION",
	"VERSION",
}

// Implement FeatureSource interface
type Source struct{}

func (s Source) Name() string { return "system" }

// NewConfig method of the FeatureSource interface
func (s *Source) NewConfig() source.Config { return nil }

// GetConfig method of the FeatureSource interface
func (s *Source) GetConfig() source.Config { return nil }

// SetConfig method of the FeatureSource interface
func (s *Source) SetConfig(source.Config) {}

func (s Source) Discover() (source.Features, error) {
	features := source.Features{}

	release, err := parseOSRelease()
	if err != nil {
		klog.Errorf("failed to get os-release: %s", err)
	} else {
		for _, key := range osReleaseFields {
			if value, exists := release[key]; exists {
				feature := "os_release." + key
				features[feature] = value

				if key == "VERSION_ID" {
					versionComponents := splitVersion(value)
					for subKey, subValue := range versionComponents {
						if subValue != "" {
							features[feature+"."+subKey] = subValue
						}
					}
				}
			}
		}
	}
	return features, nil
}

// Read and parse os-release file
func parseOSRelease() (map[string]string, error) {
	release := map[string]string{}

	f, err := os.Open(source.EtcDir.Path("os-release"))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?P<key>\w+)=(?P<value>.+)`)

	// Read line-by-line
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			release[m[1]] = strings.Trim(m[2], `"`)
		}
	}

	return release, nil
}

// Split version number into sub-components. Verifies that they are numerical
// so that they can be fully utilized in k8s nodeAffinity
func splitVersion(version string) map[string]string {
	components := map[string]string{}
	// Currently, split out just major and minor version
	re := regexp.MustCompile(`^(?P<major>\d+)(\.(?P<minor>\d+))?(\..*)?$`)
	if m := re.FindStringSubmatch(version); m != nil {
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				components[name] = m[i]
			}
		}
	}
	return components
}
