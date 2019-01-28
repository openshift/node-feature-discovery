/*
Copyright 2017 The Kubernetes Authors.

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

package rdt

import (
	"sigs.k8s.io/node-feature-discovery/source"
)

// Source implements FeatureSource.
type Source struct{}

// Name returns an identifier string for this feature source.
func (s Source) Name() string { return "rdt" }

// Discover returns feature names for CMT and CAT if supported.
func (s Source) Discover() (source.Features, error) {
	features := source.Features{}

	// Read cpuid information
	leaf07h := Cpuid(0x7, 0)
	leaf0fh := Cpuid(0xf, 0)
	leaf10h := Cpuid(0x10, 0)
	leaf07h_1 := Cpuid(0xf, 0x1)

	// Detect RDT monitoring capabilities
	if leaf07h.Ebx&(1<<12) != 0 {
		if leaf0fh.Edx&(1<<1) != 0 {
			// Monitoring is supported
			features["RDTMON"] = true

			// Cache Monitoring Technology (L3 occupancy monitoring)
			if leaf07h_1.Edx&(1<<0) != 0 {
				features["RDTCMT"] = true
			}
			// Memore Bandwidth Monitoring (L3 local&total bandwidth monitoring)
			if leaf07h_1.Edx&(3<<1) == (3 << 1) {
				features["RDTMBM"] = true
			}
		}
	}

	// Detect RDT allocation capabilities
	if leaf07h.Ebx&(1<<15) != 0 {
		// L3 Cache Allocation
		if leaf10h.Ebx&(1<<1) != 0 {
			features["RDTL3CA"] = true
		}
		// L2 Cache Allocation
		if leaf10h.Ebx&(1<<2) != 0 {
			features["RDTL2CA"] = true
		}
		// Memory Bandwidth Allocation
		if leaf10h.Ebx&(1<<3) != 0 {
			features["RDTMBA"] = true
		}
	}

	return features, nil
}
