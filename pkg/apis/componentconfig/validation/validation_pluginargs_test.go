/*
Copyright 2022 The Kubernetes Authors.

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

package validation

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/descheduler/pkg/api"
	"sigs.k8s.io/descheduler/pkg/apis/componentconfig"
)

func TestValidateRemovePodsViolatingNodeTaintsArgs(t *testing.T) {
	testCases := []struct {
		description string
		args        *componentconfig.RemovePodsViolatingNodeTaintsArgs
		expectError bool
	}{
		{
			description: "valid namespace args, no errors",
			args: &componentconfig.RemovePodsViolatingNodeTaintsArgs{
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
				},
			},
			expectError: false,
		},
		{
			description: "invalid namespaces args, expects error",
			args: &componentconfig.RemovePodsViolatingNodeTaintsArgs{
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
					Exclude: []string{"kube-system"},
				},
			},
			expectError: true,
		},
		{
			description: "valid label selector args, no errors",
			args: &componentconfig.RemovePodsViolatingNodeTaintsArgs{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"role.kubernetes.io/node": ""},
				},
			},
			expectError: false,
		},
		{
			description: "invalid label selector args, expects errors",
			args: &componentconfig.RemovePodsViolatingNodeTaintsArgs{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Operator: metav1.LabelSelectorOpIn,
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidateRemovePodsViolatingNodeTaintsArgs(tc.args)

			hasError := err != nil
			if tc.expectError != hasError {
				t.Error("unexpected arg validation behavior")
			}
		})
	}
}

func TestValidateRemovePodsViolatingNodeAffinityArgs(t *testing.T) {
	testCases := []struct {
		description string
		args        *componentconfig.RemovePodsViolatingNodeAffinityArgs
		expectError bool
	}{
		{
			description: "nil NodeAffinityType args, expects errors",
			args: &componentconfig.RemovePodsViolatingNodeAffinityArgs{
				NodeAffinityType: nil,
			},
			expectError: true,
		},
		{
			description: "empty NodeAffinityType args, expects errors",
			args: &componentconfig.RemovePodsViolatingNodeAffinityArgs{
				NodeAffinityType: []string{},
			},
			expectError: true,
		},
		{
			description: "valid NodeAffinityType args, no errors",
			args: &componentconfig.RemovePodsViolatingNodeAffinityArgs{
				NodeAffinityType: []string{"requiredDuringSchedulingIgnoredDuringExecution"},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidateRemovePodsViolatingNodeAffinityArgs(tc.args)

			hasError := err != nil
			if tc.expectError != hasError {
				t.Error("unexpected arg validation behavior")
			}
		})
	}
}

func TestValidateRemovePodLifeTimeArgs(t *testing.T) {
	testCases := []struct {
		description string
		args        *componentconfig.PodLifeTimeArgs
		expectError bool
	}{
		{
			description: "valid arg, no errors",
			args: &componentconfig.PodLifeTimeArgs{
				MaxPodLifeTimeSeconds: func(i uint) *uint { return &i }(1),
				States:                []string{string(v1.PodRunning)},
			},
			expectError: false,
		},
		{
			description: "nil MaxPodLifeTimeSeconds arg, expects errors",
			args: &componentconfig.PodLifeTimeArgs{
				MaxPodLifeTimeSeconds: nil,
			},
			expectError: true,
		},
		{
			description: "invalid pod state arg, expects errors",
			args: &componentconfig.PodLifeTimeArgs{
				States: []string{string(v1.NodeRunning)},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidatePodLifeTimeArgs(tc.args)

			hasError := err != nil
			if tc.expectError != hasError {
				t.Error("unexpected arg validation behavior")
			}
		})
	}
}
