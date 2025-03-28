/*
Copyright 2021 The Fluid Authors.
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

package jindo

import (
	"testing"

	datav1alpha1 "github.com/fluid-cloudnative/fluid/api/v1alpha1"
	"github.com/fluid-cloudnative/fluid/pkg/common"
	"github.com/fluid-cloudnative/fluid/pkg/ddc/base"
	"github.com/fluid-cloudnative/fluid/pkg/utils/fake"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newJindoEngineRT(client client.Client, name string, namespace string, withRuntimeInfo bool) *JindoEngine {
	runTimeInfo, _ := base.BuildRuntimeInfo(name, namespace, common.JindoRuntime)
	engine := &JindoEngine{
		runtime:     &datav1alpha1.JindoRuntime{},
		name:        name,
		namespace:   namespace,
		Client:      client,
		runtimeInfo: nil,
		Log:         fake.NullLogger(),
	}

	if withRuntimeInfo {
		engine.runtimeInfo = runTimeInfo
	}
	return engine
}

// TestGetRuntimeInfo tests the jindoEngine's getRuntimeInfo method under various scenarios.
// It validates the correct retrieval of runtime information from both existing and non-existing configurations,
// including error handling and nil return checks. Test cases cover:
// - Multiple JindoRuntime instances (hbase, hadoop) in the "fluid" namespace
// - Presence/Absence of runtime information (controlled by withRuntimeInfo flag)
// - Associated DaemonSets and Dataset configurations
// The test uses a fake Kubernetes client with preloaded JindoRuntimes, DaemonSets, and Datasets
// to simulate cluster state and verify expected behaviors.
func TestGetRuntimeInfo(t *testing.T) {
	runtimeInputs := []*datav1alpha1.JindoRuntime{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase",
				Namespace: "fluid",
			},
			Spec: datav1alpha1.JindoRuntimeSpec{
				Fuse: datav1alpha1.JindoFuseSpec{},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop",
				Namespace: "fluid",
			},
			Spec: datav1alpha1.JindoRuntimeSpec{
				Fuse: datav1alpha1.JindoFuseSpec{},
			},
		},
	}
	daemonSetInputs := []*v1.DaemonSet{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase-worker",
				Namespace: "fluid",
			},
			Spec: v1.DaemonSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{NodeSelector: map[string]string{"data.fluid.io/storage-fluid-hbase": "selector"}},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop-worker",
				Namespace: "fluid",
			},
			Spec: v1.DaemonSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{NodeSelector: map[string]string{"data.fluid.io/storage-fluid-hadoop": "selector"}},
				},
			},
		},
	}
	dataSetInputs := []*datav1alpha1.Dataset{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop",
				Namespace: "fluid",
			},
		},
	}
	objs := []runtime.Object{}
	for _, runtimeInput := range runtimeInputs {
		objs = append(objs, runtimeInput.DeepCopy())
	}
	for _, daemonSetInput := range daemonSetInputs {
		objs = append(objs, daemonSetInput.DeepCopy())
	}
	for _, dataSetInput := range dataSetInputs {
		objs = append(objs, dataSetInput.DeepCopy())
	}
	//scheme := runtime.NewScheme()
	//scheme.AddKnownTypes(v1.SchemeGroupVersion, daemonSetWithSelector)
	//scheme.AddKnownTypes(v1alpha1.GroupVersion,runtimeInput)
	fakeClient := fake.NewFakeClientWithScheme(testScheme, objs...)

	testCases := []struct {
		name            string
		namespace       string
		withRuntimeInfo bool
		isErr           bool
		isNil           bool
	}{
		{
			name:            "hbase",
			namespace:       "fluid",
			withRuntimeInfo: false,
			isErr:           false,
			isNil:           false,
		},
		{
			name:            "hbase",
			namespace:       "fluid",
			withRuntimeInfo: false,
			isErr:           false,
			isNil:           false,
		},
		{
			name:            "hbase",
			namespace:       "fluid",
			withRuntimeInfo: true,
			isErr:           false,
			isNil:           false,
		},
		{
			name:            "hadoop",
			namespace:       "fluid",
			withRuntimeInfo: false,
			isErr:           false,
			isNil:           false,
		},
	}
	for _, testCase := range testCases {
		engine := newJindoEngineRT(fakeClient, testCase.name, testCase.namespace, testCase.withRuntimeInfo)
		runtimeInfo, err := engine.getRuntimeInfo()
		isNil := runtimeInfo == nil
		isErr := err != nil
		if isNil != testCase.isNil {
			t.Errorf(" want %t, got %t", testCase.isNil, isNil)
		}
		if isErr != testCase.isErr {
			t.Errorf(" want %t, got %t", testCase.isErr, isErr)
		}
	}
}
