/*
Copyright 2019 The Kubernetes Authors.

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

package main

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vektra/errors"
	api "k8s.io/api/core/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/node-feature-discovery/pkg/apihelper"
	"sigs.k8s.io/node-feature-discovery/pkg/version"
)

func TestDiscoveryWithMockSources(t *testing.T) {
	Convey("When I discover features from fake source and update the node using fake client", t, func() {
		fakeFeatureLabels := map[string]string{"source-feature.1": "val1", "source-feature.2": "val2", "source-feature.3": "val3"}
		fakeAnnotations := map[string]string{"version": version.Get()}
		fakeFeatureLabelNames := make([]string, 0, len(fakeFeatureLabels))
		for k, _ := range fakeFeatureLabels {
			fakeFeatureLabelNames = append(fakeFeatureLabelNames, k)
		}
		fakeAnnotations["feature-labels"] = strings.Join(fakeFeatureLabelNames, ",")

		mockAPIHelper := new(apihelper.MockAPIHelpers)
		testHelper := apihelper.APIHelpers(mockAPIHelper)
		mockNode := &api.Node{}
		mockNodeName := "mock-node"
		var mockClient *k8sclient.Clientset

		Convey("When I successfully update the node with feature labels", func() {
			mockAPIHelper.On("GetClient").Return(mockClient, nil)
			mockAPIHelper.On("GetNode", mockClient, mockNodeName).Return(mockNode, nil).Once()
			mockAPIHelper.On("AddLabels", mockNode, fakeFeatureLabels).Return().Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, labelNs).Return().Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, "node.alpha.kubernetes-incubator.io/nfd").Return().Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, "node.alpha.kubernetes-incubator.io/node-feature-discovery").Return().Once()
			mockAPIHelper.On("AddAnnotations", mockNode, fakeAnnotations).Return().Once()
			mockAPIHelper.On("UpdateNode", mockClient, mockNode).Return(nil).Once()
			err := updateNodeFeatures(testHelper, mockNodeName, fakeFeatureLabels, fakeAnnotations)

			Convey("Error is nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I fail to update the node with feature labels", func() {
			expectedError := errors.New("fake error")
			mockAPIHelper.On("GetClient").Return(nil, expectedError)
			err := updateNodeFeatures(testHelper, mockNodeName, fakeFeatureLabels, fakeAnnotations)

			Convey("Error is produced", func() {
				So(err, ShouldEqual, expectedError)
			})
		})

		Convey("When I fail to get a mock client while updating feature labels", func() {
			expectedError := errors.New("fake error")
			mockAPIHelper.On("GetClient").Return(nil, expectedError)
			err := updateNodeFeatures(testHelper, mockNodeName, fakeFeatureLabels, fakeAnnotations)

			Convey("Error is produced", func() {
				So(err, ShouldEqual, expectedError)
			})
		})

		Convey("When I fail to get a mock node while updating feature labels", func() {
			expectedError := errors.New("fake error")
			mockAPIHelper.On("GetClient").Return(mockClient, nil)
			mockAPIHelper.On("GetNode", mockClient, mockNodeName).Return(nil, expectedError).Once()
			err := updateNodeFeatures(testHelper, mockNodeName, fakeFeatureLabels, fakeAnnotations)

			Convey("Error is produced", func() {
				So(err, ShouldEqual, expectedError)
			})
		})

		Convey("When I fail to update a mock node while updating feature labels", func() {
			expectedError := errors.New("fake error")
			mockAPIHelper.On("GetClient").Return(mockClient, nil)
			mockAPIHelper.On("GetNode", mockClient, mockNodeName).Return(mockNode, nil).Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, labelNs).Return().Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, "node.alpha.kubernetes-incubator.io/nfd").Return().Once()
			mockAPIHelper.On("RemoveLabelsWithPrefix", mockNode, "node.alpha.kubernetes-incubator.io/node-feature-discovery").Return().Once()
			mockAPIHelper.On("AddLabels", mockNode, fakeFeatureLabels).Return().Once()
			mockAPIHelper.On("AddAnnotations", mockNode, fakeAnnotations).Return().Once()
			mockAPIHelper.On("UpdateNode", mockClient, mockNode).Return(expectedError).Once()
			err := updateNodeFeatures(testHelper, mockNodeName, fakeFeatureLabels, fakeAnnotations)

			Convey("Error is produced", func() {
				So(err, ShouldEqual, expectedError)
			})
		})

	})
}

func TestArgsParse(t *testing.T) {
	Convey("When parsing command line arguments", t, func() {
		argv1 := []string{"--no-publish"}
		argv2 := []string{"--label-whitelist=.*rdt.*"}

		Convey("When --no-publish and --oneshot flags are passed", func() {
			args := argsParse(argv1)

			Convey("noPublish is set and args.sources is set to the default value", func() {
				So(args.noPublish, ShouldBeTrue)
				So(len(args.labelWhiteList.String()), ShouldEqual, 0)
			})
		})

		Convey("When --label-whitelist flag is passed and set to some value", func() {
			args := argsParse(argv2)

			Convey("args.labelWhiteList is set to appropriate value and args.sources is set to default value", func() {
				So(args.noPublish, ShouldBeFalse)
				So(args.labelWhiteList.String(), ShouldResemble, ".*rdt.*")
			})
		})
	})
}
