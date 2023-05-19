// Copyright 2022 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/types"
	testutils_test "sigs.k8s.io/kustomize/kustomize/v5/commands/internal/testutils"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func TestAddBuildMetadata(t *testing.T) {
	tests := map[string]struct {
		input       string
		args        []string
		expectedErr string
	}{
		"happy path": {
			input: ``,
			args:  []string{strings.Join(types.BuildMetadataOptions, ",")},
		},
		"option already there": {
			input: `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
buildMetadata: [originAnnotations]`,
			args:        []string{types.OriginAnnotations},
			expectedErr: "buildMetadata option originAnnotations already in kustomization file",
		},
		"invalid option": {
			input:       ``,
			args:        []string{"invalid_option"},
			expectedErr: "invalid buildMetadata option: invalid_option",
		},
		"too many args": {
			input:       ``,
			args:        []string{"option1", "option2"},
			expectedErr: "too many arguments: [option1 option2]; to provide multiple buildMetadata options, please separate options by comma",
		},
	}

	for _, tc := range tests {
		fSys := filesys.MakeFsInMemory()
		testutils_test.WriteTestKustomizationWith(fSys, []byte(tc.input))
		cmd := newCmdAddBuildMetadata(fSys)
		err := cmd.RunE(cmd, tc.args)
		if tc.expectedErr != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		} else {
			assert.NoError(t, err)
			content, err := testutils_test.ReadTestKustomization(fSys)
			assert.NoError(t, err)
			for _, opt := range strings.Split(tc.args[0], ",") {
				assert.Contains(t, string(content), opt)
			}
		}
	}
}
