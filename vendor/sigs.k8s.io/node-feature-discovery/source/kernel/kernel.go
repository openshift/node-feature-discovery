/*
Copyright 2018 The Kubernetes Authors.

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

package kernel

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"sigs.k8s.io/node-feature-discovery/source"
)

// Configuration file options
type NFDConfig struct {
	KconfigFile string
	ConfigOpts  []string `json:"configOpts,omitempty"`
}

var Config = NFDConfig{
	KconfigFile: "",
	ConfigOpts: []string{
		"NO_HZ",
		"NO_HZ_IDLE",
		"NO_HZ_FULL",
		"PREEMPT",
	},
}

// Implement FeatureSource interface
type Source struct{}

func (s Source) Name() string { return "kernel" }

func (s Source) Discover() (source.Features, error) {
	features := source.Features{}

	// Read kernel version
	version, err := parseVersion()
	if err != nil {
		log.Printf("ERROR: Failed to get kernel version: %s", err)
	} else {
		for key := range version {
			features["version."+key] = version[key]
		}
	}

	// Read kconfig
	kconfig, err := parseKconfig()
	if err != nil {
		log.Printf("ERROR: Failed to read kconfig: %s", err)
	}

	// Check flags
	for _, opt := range Config.ConfigOpts {
		if _, ok := kconfig[opt]; ok {
			features["config."+opt] = true
		}
	}

	selinux, err := SelinuxEnabled()
	if err != nil {
		log.Print(err)
	} else if selinux {
		features["selinux.enabled"] = true
	}

	return features, nil
}

// Read and parse kernel version
func parseVersion() (map[string]string, error) {
	version := map[string]string{}

	// Open file for reading
	raw, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return nil, err
	}

	full := strings.TrimSpace(string(raw))
	version["full"] = full

	// Regexp for parsing version components
	re := regexp.MustCompile(`^(?P<major>\d+)(\.(?P<minor>\d+))?(\.(?P<revision>\d+))?(-.*)?$`)
	if m := re.FindStringSubmatch(full); m != nil {
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				version[name] = m[i]
			}
		}
	}

	return version, nil
}

// Read gzipped kernel config
func readKconfigGzip(filename string) ([]byte, error) {
	// Open file for reading
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Uncompress data
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

// Read kconfig into a map
func parseKconfig() (map[string]bool, error) {
	kconfig := map[string]bool{}
	raw := []byte(nil)
	err := error(nil)

	// First, try kconfig specified in the config file
	if len(Config.KconfigFile) > 0 {
		raw, err = ioutil.ReadFile(Config.KconfigFile)
		if err != nil {
			log.Printf("ERROR: Failed to read kernel config from %s: %s", Config.KconfigFile, err)
		}
	}

	// Then, try to read from /proc
	if raw == nil {
		raw, err = readKconfigGzip("/proc/config.gz")
		if err != nil {
			log.Printf("Failed to read /proc/config.gz: %s", err)
		}
	}

	// Last, try to read from /boot/
	if raw == nil {
		// Get kernel version
		unameRaw, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
		uname := strings.TrimSpace(string(unameRaw))
		if err != nil {
			return nil, err
		}
		// Read kconfig
		raw, err = ioutil.ReadFile("/host-boot/config-" + uname)
		if err != nil {
			return nil, err
		}
	}

	// Regexp for matching kconfig flags
	re := regexp.MustCompile(`^CONFIG_(?P<flag>\w+)=(?P<value>.+)`)

	// Process data, line-by-line
	lines := bytes.Split(raw, []byte("\n"))
	for _, line := range lines {
		if m := re.FindStringSubmatch(string(line)); m != nil {
			if m[2] == "y" || m[2] == "m" {
				kconfig[m[1]] = true
			}
		}
	}

	return kconfig, nil
}
