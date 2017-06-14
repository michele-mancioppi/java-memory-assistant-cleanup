/*
 * Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file at the root of the repository.
 */

package matchers

import (
	"fmt"
	"os"

	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

// HaveFile checks if the afero FS has a given file in it
func HaveFile(expected interface{}) types.GomegaMatcher {
	return &hasFile{
		expected: expected,
	}
}

type hasFile struct {
	expected interface{}
}

func (matcher *hasFile) Match(actual interface{}) (success bool, err error) {
	fs, ok := actual.(afero.Fs)
	if !ok {
		return false, fmt.Errorf("HaveFile matcher expects an afero.Fs as 'actual'")
	}

	fileName, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveFile matcher expects a string as 'expected'")
	}

	if _, err := fs.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("Cannot open file '%v': %s", fileName, err.Error())
	}

	return true, nil
}

func (matcher *hasFile) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain a file named\n\t%#v", actual, matcher.expected)
}

func (matcher *hasFile) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain a file named\n\t%#v", actual, matcher.expected)
}
