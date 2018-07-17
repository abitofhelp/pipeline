////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package cmd implements the command-line actions for the application.
package cmd

import (
	"errors"
	"fmt"
	. "github.com/abitofhelp/go-helpers/error"
	. "github.com/abitofhelp/pipeline/pipeline"
	"gopkg.in/urfave/cli.v2"
)

const (
	// The default number of bytes for a reusable buffer that is used when scanning for files.
	kDefaultScannerBufferSize = 64 * 1024

	// The default number of paths that can be loaded into the paths channel.
	kDefaultPathsChannelSize = 100

	// The default number of goroutines that will consume the paths channel.
	kDefaultPathConsumerCount = 20

	// The maximum scanner buffer size is 1GB.
	kMaxScannerBufferSize = 1000 * 1024

	// The maximum number of paths that can be loaded into the paths channel.
	kMaxPathsChannelSize = 256

	// The maximum number of goroutine consumers on the paths channel.
	kMaxPathConsumerCount = 50
)

// Variable Start is a command that defines the command line for starting the pipeline's processing.
var (
	Start = cli.Command{
		Name:   "start",
		Usage:  "Starts the pipeline's processing",
		Action: start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "path to the directory containing files to process",
				Value: "/tmp",
			},
			&cli.Uint64Flag{
				Name:  "sbs",
				Usage: "(scannerBufferSize) is the number of reusable bytes to use for the directory scanner's work",
				Value: kDefaultScannerBufferSize, // allocate once and re-use each time - faster...
			},
			&cli.Uint64Flag{
				Name:  "pcs",
				Usage: "(pathChanSize) is the number of file system paths that will be buffered in a channel in the pipeline",
				Value: kDefaultPathsChannelSize,
			},
			&cli.Uint64Flag{
				Name:  "pcc",
				Usage: "(pathConsumerCount) is the number of concurrent and parallel goroutines that will consume paths from a channel in the pipeline",
				Value: kDefaultPathConsumerCount,
			},
		},
	}

	// Variable APipeline is the pipeline that has been created using the command-line information.
	// It is a global value that others can use.
	APipeline IPipeline = nil
)

// Function start is an internal function that validates the start command's
// parameters and completes other configuration activities.
// Returns nil is there are no errors.
func start(c *cli.Context) (err error) {

	// Validate that the command line is correct for processing.
	err = validateCommandLine(c)
	if IsError(err, nil) {
		return err
	}

	// Create the pipeline and Start it.
	APipeline, err = createPipeline(c)
	if IsError(err, nil) {
		return err
	}

	// Start the pipeline...
	err = APipeline.Start()
	if IsError(err, nil) {
		return err
	}

	return nil
}

// Function buildPipeline uses the command-line parameters to create a properly
// configured pipeline.
func createPipeline(c *cli.Context) (IPipeline, error) {
	var (
		path              = c.String("path")
		scannerBufferSize = c.Uint64("sbs")
		pathChanSize      = c.Uint64("pcs")
		pathConsumerCount = c.Uint64("pcc")
	)

	// Create an instance of the pipeline using our command-line options.
	pipeline, err := New(path, scannerBufferSize, pathChanSize, pathConsumerCount)
	if IsError(err, nil) {
		return nil, err
	}

	return pipeline, nil
}

// Function validateCommandLine is an internal function that examines
// the command line parameters to ensure that they are correct.
func validateCommandLine(c *cli.Context) (err error) {
	var (
		path              = c.String("path")
		scannerBufferSize = c.Uint64("sbs")
		pathChanSize      = c.Uint64("pcs")
		pathConsumerCount = c.Uint64("pcc")
	)

	err = nil

	switch {

	case path == "":
		err = errors.New("there must be a path to the directory containing files to process")

	case scannerBufferSize > kMaxScannerBufferSize:
		err = errors.New(fmt.Sprintf("%s%d bytes", "the maximum scanner buffer size cannot exceed", kMaxScannerBufferSize))

	case pathChanSize > kMaxPathsChannelSize:
		err = errors.New(fmt.Sprintf("%s%d bytes", "the maximum number of paths in the channel cannnot exceed", kMaxPathsChannelSize))

	case pathConsumerCount > kMaxPathConsumerCount:
		err = errors.New(fmt.Sprintf("%s%d bytes", "the maximum number of goroutine consumers on the paths channel cannot exceed", kMaxPathConsumerCount))
	}

	return err
}
