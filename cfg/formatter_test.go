// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_formatGoFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(mfs *mockFileSystem, mf *mockFormatter)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(mfs *mockFileSystem, mf *mockFormatter) {},
		},
		{
			name: "error when reading file",
			mockClosure: func(mfs *mockFileSystem, mf *mockFormatter) {
				mfs.readFileErr = errors.New("read error")
			},
			expectedError: errors.New("reading go file file.go: read error"),
		},
		{
			name: "error when formatting file",
			mockClosure: func(mfs *mockFileSystem, mf *mockFormatter) {
				mf.sourceErr = errors.New("source error")
			},
			expectedError: errors.New("formating go file file.go: source error"),
		},
		{
			name: "error when writing file",
			mockClosure: func(mfs *mockFileSystem, mf *mockFormatter) {
				mfs.writeFileErr = errors.New("write error")
			},
			expectedError: errors.New("writing go file file.go: write error"),
		},
	}
	for _, tc := range testCases {
		mfs := new(mockFileSystem)
		mf := new(mockFormatter)
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfs, mf)
			fsProvider = mfs
			formatterProvider = mf
			err := formatGoFile("file.go")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}
