/*
 * Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file at the root of the repository.
 */

package main_test

import (
	"os"

	"github.com/spf13/afero"

	. "github.com/SAP/java-memory-assistant/cleanup"
	. "github.com/SAP/java-memory-assistant/cleanup/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run clean_up", func() {

	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
	})

	Context("without 'JMA_HEAP_DUMP_FOLDER' environment variable set", func() {

		It("fails", func() {
			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount: 1,
			})

			Expect(err.Error()).To(Equal("The environment variable 'JMA_HEAP_DUMP_FOLDER' is not set"))
			Expect(deletedFiles).To(BeEmpty())
		})

	})

	Context("with the 'JMA_HEAP_DUMP_FOLDER' pointing to a non-existing folder", func() {

		It("fails", func() {
			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount:   1,
				HeapDumpFolder: "nope",
			})

			Expect(err.Error()).To(ContainSubstring("Cannot open 'JMA_HEAP_DUMP_FOLDER' directory 'nope': does not exist"))
			Expect(deletedFiles).To(BeEmpty())
		})

	})

	Context("with the 'JMA_HEAP_DUMP_FOLDER' pointing to a regular file", func() {

		It("fails", func() {
			fs.Create("dumps")

			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount:   1,
				HeapDumpFolder: "dumps",
			})

			Expect(err.Error()).To(Equal("Cannot open 'JMA_HEAP_DUMP_FOLDER' directory 'dumps': not a directory (mode: T---------)"))
			Expect(deletedFiles).To(BeEmpty())
		})

	})

	Context("with 3 heap dump files", func() {

		BeforeEach(func() {
			fs.MkdirAll("dumps", os.ModeDir)
			fs.Create("dumps/1.hprof")
			fs.Create("dumps/2.hprof")
			fs.Create("dumps/3.hprof")
			fs.Create("dumps/not.a.dump")
		})

		AfterEach(func() {
			if _, err := fs.Stat("dumps/not.a.dump"); os.IsNotExist(err) {
				Fail("A non-dump file has been deleted")
			}
		})

		It("with max one heap dump, it deletes all three files", func() {
			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount:   1,
				HeapDumpFolder: "dumps"})

			Expect(err).To(BeNil())
			Expect(deletedFiles).To(ConsistOf("dumps/1.hprof", "dumps/2.hprof", "dumps/3.hprof"))

			Expect(fs).ToNot(HaveFile("dumps/1.hprof"))
			Expect(fs).ToNot(HaveFile("dumps/2.hprof"))
			Expect(fs).ToNot(HaveFile("dumps/3.hprof"))
		})

		It("with max two heap dump, it deletes the first two files", func() {
			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount:   2,
				HeapDumpFolder: "dumps"})

			Expect(err).To(BeNil())
			Expect(deletedFiles).To(ConsistOf("dumps/1.hprof", "dumps/2.hprof"))

			Expect(fs).ToNot(HaveFile("dumps/1.hprof"))
			Expect(fs).ToNot(HaveFile("dumps/2.hprof"))
			Expect(fs).To(HaveFile("dumps/3.hprof"))
		})
	})

	Context("with repeated invocations", func() {

		BeforeEach(func() {
			fs.MkdirAll("dumps", os.ModeDir)
			fs.Create("dumps/1.hprof")
			fs.Create("dumps/2.hprof")
			fs.Create("dumps/not.a.dump")
		})

		AfterEach(func() {
			if _, err := fs.Stat("dumps/not.a.dump"); os.IsNotExist(err) {
				Fail("A non-dump file has been deleted")
			}
		})

		It("is idempotent", func() {
			deletedFiles_1, err_1 := CleanUp(fs, Config{
				MaxDumpCount:   2,
				HeapDumpFolder: "dumps",
			})

			Expect(err_1).To(BeNil())
			Expect(deletedFiles_1).To(ConsistOf("dumps/1.hprof"))

			Expect(fs).ToNot(HaveFile("dumps/1.hprof"))
			Expect(fs).To(HaveFile("dumps/2.hprof"))

			// Test repeated invocations
			deletedFiles_2, err_2 := CleanUp(fs, Config{
				MaxDumpCount:   2,
				HeapDumpFolder: "dumps",
			})

			Expect(err_2).To(BeNil())
			Expect(deletedFiles_2).To(BeEmpty())

			Expect(fs).To(HaveFile("dumps/2.hprof"))
		})

	})

	Context("with no heap dump files", func() {

		BeforeEach(func() {
			fs.MkdirAll("dumps", os.ModeDir)
			fs.Create("dumps/not.a.dump")
		})

		AfterEach(func() {
			if _, err := fs.Stat("dumps/not.a.dump"); os.IsNotExist(err) {
				Fail("A non-dump file has been deleted")
			}
		})

		It("with max one heap dump, it deletes no files", func() {
			deletedFiles, err := CleanUp(fs, Config{
				MaxDumpCount:   3,
				HeapDumpFolder: "dumps"})

			Expect(err).To(BeNil())
			Expect(deletedFiles).To(BeEmpty())
		})
	})

})
