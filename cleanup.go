/*
 * Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file at the root of the repository.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"robpike.io/filter"

	"github.com/caarlos0/env"
	"github.com/spf13/afero"
)

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	deletedFiles, err := CleanUp(afero.NewOsFs(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, deletedFile := range deletedFiles {
		fmt.Printf("Heap dump '%v' deleted\n", deletedFile)
	}
}

type byName []os.FileInfo

// Config visible for testing
type Config struct {
	HeapDumpFolder string `env:"JMA_HEAP_DUMP_FOLDER"`
	MaxDumpCount   int    `env:"JMA_MAX_DUMP_COUNT" envDefault:"0"`
}

func (f byName) Len() int {
	return len(f)
}

func (f byName) Less(i, j int) bool {
	return f[i].Name() < f[j].Name()
}

func (f byName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// CleanUp visible for testing
func CleanUp(fs afero.Fs, cfg Config) ([]string, error) {
	if cfg.HeapDumpFolder == "" {
		return nil, fmt.Errorf("The environment variable 'JMA_HEAP_DUMP_FOLDER' is not set")
	}

	heapDumpFolder := cfg.HeapDumpFolder

	maxDumpCount := cfg.MaxDumpCount
	if maxDumpCount < 0 {
		return nil, fmt.Errorf("The value of the 'JMA_MAX_DUMP_COUNT' environment variable contains a negative number: %v", maxDumpCount)
	}

	file, err := fs.Stat(heapDumpFolder)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Cannot open 'JMA_HEAP_DUMP_FOLDER' directory '%v': does not exist", heapDumpFolder)
	}

	mode := file.Mode()
	if !mode.IsDir() {
		return nil, fmt.Errorf("Cannot open 'JMA_HEAP_DUMP_FOLDER' directory '%v': not a directory (mode: %v)", heapDumpFolder, mode)
	}

	files, err := afero.ReadDir(fs, heapDumpFolder)
	if err != nil {
		return nil, fmt.Errorf("Cannot open 'JMA_HEAP_DUMP_FOLDER' directory '%v': %v", heapDumpFolder, err)
	}

	isHeapDumpFile := func(file os.FileInfo) bool {
		return strings.HasSuffix(file.Name(), ".hprof")
	}

	heapDumpFiles := filter.Choose(files, isHeapDumpFile).([]os.FileInfo)

	if len(heapDumpFiles) < maxDumpCount || maxDumpCount < 1 {
		return []string{}, nil
	}

	var deletedFiles []string
	sort.Sort(sort.Reverse(byName(heapDumpFiles)))

	for _, file := range heapDumpFiles[maxDumpCount-1:] {
		path := heapDumpFolder + "/" + file.Name()
		var err = fs.Remove(path)
		if err != nil {
			return nil, fmt.Errorf("Cannot delete heap dump file '"+path+"': %v", err)
		}

		deletedFiles = append(deletedFiles, path)
	}

	return deletedFiles, nil
}
