/*
 * Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file at the root of the repository.
 */

package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCleanUp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JavaMemoryAssistant CleanUp Suite")
}
