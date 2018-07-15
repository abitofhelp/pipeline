////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a BSD-style  license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package pipeline implements a processing pipeline with multiple steps.
package pipeline

import (
	"fmt"
	godirwalk "github.com/karrick/godirwalk"
	"os"
	. "path/filepath"
	"sync"
)

// Recursively walk a file system hierarchy to locate files, and pass the paths into the pipeline for processing.
// Parameter pathToDirectory is the path to a folder containing files to process.
// Parameter ch is the unidirectional channel being used to feed the paths to the pipeline.
func (p Pipeline) loadPathsToChannel(pathToDirectory string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	cnt := uint64(0)

	godirwalk.Walk(pathToDirectory, &godirwalk.Options{

		FollowSymbolicLinks: false,
		Unsorted:            true,

		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				wg.Add(1)
				fmt.Printf("\ncnt: %d", cnt)
				ch <- Join(path, de.Name())
				cnt++
			}

			// Signal no errors...
			return nil
		},

		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			// Your program may want to log the error somehow.
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)

			// On error, we will skip the current file system node and continue
			// walking the file system hierarchy of remaining nodes.
			return godirwalk.SkipNode
		},
	})
}
