////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a BSD-style  license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package main contains the entry point to the application and configuration settings.
package main

import (
	"fmt"
	. "github.com/abitofhelp/pipeline/pipeline"
	"os"
	"path/filepath"
	"sync"
)

// Function main is the entry point to the application.
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s directory]]\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	// Create the pipeline...

	pipeline, err := New(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create the pipeline: %s", err)
		os.Exit(2)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Construct the pipeline...
	go pipeline.Build(&wg)

	// Start processing by the pipeline...

	// Abort processing by the pipeline...  Handle exit gracefully.

	// End processing by the pipeline...

	// Wait for goroutines to finish...
	wg.Wait()

}
