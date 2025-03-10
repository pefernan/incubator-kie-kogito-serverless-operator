/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package platform

import (
	"os"
	"regexp"
	"testing"

	"github.com/apache/incubator-kie-kogito-serverless-operator/api/v1alpha08"
	"github.com/apache/incubator-kie-kogito-serverless-operator/controllers/cfg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"

	"github.com/apache/incubator-kie-kogito-serverless-operator/test"
)

const dockerFile = "FROM docker.io/apache/default-test-kie-sonataflow-builder:main AS builder\n\n# ETC, \n\n# ETC, \n\n# ETC"

func TestSonataFlowBuildController(t *testing.T) {
	platform := test.GetBasePlatform()
	dockerfileBytes, err := os.ReadFile("../../test/builder/Dockerfile")
	if err != nil {
		assert.Fail(t, "Unable to read base Dockerfile")
	}
	dockerfile := string(dockerfileBytes)
	// 1 - Let's verify that the default image is used (for this unit test is docker.io/apache/incubator-kie-sonataflow-builder:main)
	resDefault := GetCustomizedBuilderDockerfile(dockerfile, *platform)
	foundDefault, err := regexp.MatchString("FROM docker.io/apache/incubator-kie-sonataflow-builder:main AS builder", resDefault)
	assert.NoError(t, err)
	assert.True(t, foundDefault)

	// 2 - Let's try to override using the productized image
	platform.Spec.Build.Config.BaseImage = "registry.access.redhat.com/openshift-serverless-1-tech-preview/logic-swf-builder-rhel8"
	resProductized := GetCustomizedBuilderDockerfile(dockerfile, *platform)
	foundProductized, err := regexp.MatchString("FROM registry.access.redhat.com/openshift-serverless-1-tech-preview/logic-swf-builder-rhel8 AS builder", resProductized)
	assert.NoError(t, err)
	assert.True(t, foundProductized)
}

func TestGetCustomizedBuilderDockerfile_NoBaseImageCustomization(t *testing.T) {
	sfp := v1alpha08.SonataFlowPlatform{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1alpha08.SonataFlowPlatformSpec{},
		Status:     v1alpha08.SonataFlowPlatformStatus{},
	}
	customizedDockerfile := GetCustomizedBuilderDockerfile(dockerFile, sfp)
	assert.Equal(t, dockerFile, customizedDockerfile)
}

func TestGetCustomizedBuilderDockerfile_BaseImageCustomizationFromPlatform(t *testing.T) {
	sfp := v1alpha08.SonataFlowPlatform{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1alpha08.SonataFlowPlatformSpec{
			Build: v1alpha08.BuildPlatformSpec{
				Template: v1alpha08.BuildTemplate{},
				Config: v1alpha08.BuildPlatformConfig{
					BaseImage: "docker.io/apache/platfom-sonataflow-builder:main",
				},
			},
		},
		Status: v1alpha08.SonataFlowPlatformStatus{},
	}

	expectedDockerFile := "FROM docker.io/apache/platfom-sonataflow-builder:main AS builder\n\n# ETC, \n\n# ETC, \n\n# ETC"
	customizedDockerfile := GetCustomizedBuilderDockerfile(dockerFile, sfp)
	assert.Equal(t, expectedDockerFile, customizedDockerfile)
}

func TestGetCustomizedBuilderDockerfile_BaseImageCustomizationFromControllersConfig(t *testing.T) {
	sfp := v1alpha08.SonataFlowPlatform{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1alpha08.SonataFlowPlatformSpec{},
		Status:     v1alpha08.SonataFlowPlatformStatus{},
	}

	_, err := cfg.InitializeControllersCfgAt("../cfg/testdata/controllers-cfg-test.yaml")
	assert.NoError(t, err)
	expectedDockerFile := "FROM local/sonataflow-builder:1.0.0 AS builder\n\n# ETC, \n\n# ETC, \n\n# ETC"
	customizedDockerfile := GetCustomizedBuilderDockerfile(dockerFile, sfp)
	assert.Equal(t, expectedDockerFile, customizedDockerfile)
}
