////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a BSD-style  license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package pipeline implements a processing pipeline with multiple steps.
package pipeline

import (
	"errors"
	"fmt"
	. "github.com/abitofhelp/go-helpers/string"
	. "github.com/abitofhelp/go-helpers/time"
	"strings"
	"sync"
	"time"
)

// Type Status indicates the current state of the pipeline.
type Status int

// Constants for the states of a Status.
const (
	Aborting Status = iota
	Starting
	Running
	Stopping
	Stopped
)

// Type IPipeline is an interface that requires implementations of its pipeline methods.
type IPipeline interface {

	// Function Start initiates processing in the pipeline.
	Start() error

	// Function Abort abends processing in the pipeline.
	Abort() error

	// Function Stop terminates processing after all steps have been completed.
	Stop() error
}

// Type Pipeline is a struct that provides data and methods to create and manage a pipeline.
type Pipeline struct {
	// Field startedUtc is the date/time in UTC when the pipeline commences its work.
	startedUtc time.Time

	// Field endedUtc is the date/time in UTC when the pipeline completed its work.
	endedUtc time.Time

	// Field status indicates the current status of the pipeline.
	status Status

	// Field path is the file system path to a directory containing image files to process.
	path string

	// Field scannerBufferSize is  the number of reusable bytes to use for the directory scanner's work.
	scannerBufferSize uint64

	// Field pathChanSize is the number of file system paths that will be buffered in a channel in the pipeline.
	pathChanSize uint64

	// Field pathConsumerCount is the number of concurrent and parallel goroutines that will consume paths from a channel in the pipeline.
	pathConsumerCount uint64

	// Field pathsChannel is the channel containing paths to the files that will be processed.
	pathsChannel chan string

	// Field commandChannel is the channel that will signal to start the pipeline.
	commandChannel chan bool
}

// Function New is a factory that creates an initialized Pipeline.
// Parameter path to the directory containing files to process.
// Parameter scannerBufferSize is  the number of reusable bytes to use for the directory scanner's work.
// Parameter pathChanSize is the number of file system paths that will be buffered in a channel in the pipeline.
// Parameter pathConsumerCount is the number of concurrent and parallel goroutines that will consume paths from a channel in the pipeline.
// Returns an initialized pipeline or error.
func New(path string, scannerBufferSize uint64, pathChanSize uint64, pathConsumerCount uint64) (*Pipeline, error) {
	pipeline := &Pipeline{}
	if pipeline == nil {
		return nil, errors.New("failed to create an instance of Pipeline")
	}

	err := pipeline.setPath(path)
	if err != nil {
		return nil, err
	}

	err = pipeline.setScannerBufferSize(scannerBufferSize)
	if err != nil {
		return nil, err
	}

	err = pipeline.setPathChanSize(pathChanSize)
	if err != nil {
		return nil, err
	}

	err = pipeline.setPathConsumerCount(pathConsumerCount)
	if err != nil {
		return nil, err
	}

	err = pipeline.setEndedUtc(Zero())
	if err != nil {
		return nil, err
	}

	err = pipeline.setStatus(Stopped)
	if err != nil {
		return nil, err
	}

	// Create the channel that will provide paths to files for processing.
	err = pipeline.setPathsChannel(make(chan string, pipeline.PathChanSize()))
	if err != nil {
		return nil, err
	}

	// Create the the channel that will signal to start the pipeline.
	err = pipeline.setCommandChannel(make(chan bool))
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

// Method Path gets the path from the instance.
func (p Pipeline) Path() string {
	return p.path
}

// Method SetPath sets the value of the path in the instance.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setPath(path string) error {

	if path == "" {
		return errors.New("the path cannot be empty")
	}
	path = CleanStringForPlatform(path)
	p.path = path

	return nil
}

// Method ScannerBufferSize gets the number of reusable bytes to use for the directory scanner's work.
func (p Pipeline) ScannerBufferSize() uint64 {
	return p.scannerBufferSize
}

// Method setScannerBufferSize sets the number of reusable bytes to use for the directory scanner's work.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setScannerBufferSize(scannerBufferSize uint64) error {
	p.scannerBufferSize = scannerBufferSize
	return nil
}

// Method PathChanSize gets the number of file system paths that will be buffered in a channel in the pipeline.
func (p Pipeline) PathChanSize() uint64 {
	return p.pathChanSize
}

// Method setPathChanSize sets the number of file system paths that will be buffered in a channel in the pipeline.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setPathChanSize(pathChanSize uint64) error {
	p.pathChanSize = pathChanSize
	return nil
}

// Method PathConsumerCount gets the number of concurrent and parallel goroutines that will consume paths from a channel in the pipeline.
func (p Pipeline) PathConsumerCount() uint64 {
	return p.pathConsumerCount
}

// Method setPathConsumerCount sets the number of concurrent and parallel goroutines that will consume paths from a channel in the pipeline.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setPathConsumerCount(pathConsumerCount uint64) error {
	p.pathConsumerCount = pathConsumerCount
	return nil
}

// Method Status gets the current status from the instance of a Pipeline.
func (p Pipeline) Status() Status {
	return p.status
}

// Method SetStatus sets the status of the Pipeline.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setStatus(status Status) error {
	p.status = status
	return nil
}

// Method Started gets the UTC date/time when the pipeline started processing.
func (p Pipeline) StartedUtc() time.Time {
	return p.startedUtc
}

// Method SetStartedUtc sets the value of the startedUtc, when the pipeline started processing.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setStartedUtc(startedUtc time.Time) error {

	utc := CleanStringForPlatform(startedUtc.Location().String())
	if strings.Compare(utc, "UTC") != 0 {
		return errors.New("the startedUtc value must be in UTC")
	}

	p.startedUtc = startedUtc

	return nil
}

// Method Ended gets the UTC date/time when the pipeline completed processing.
func (p Pipeline) EndedUtc() time.Time {
	return p.endedUtc
}

// Method SetEndedUtc sets the value of the endedUtc, when the pipeline completed processing.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setEndedUtc(endedUtc time.Time) error {

	utc := CleanStringForPlatform(endedUtc.Location().String())
	if strings.Compare(utc, "UTC") != 0 {
		return errors.New("the endedUtc value must be in UTC")
	}

	p.endedUtc = endedUtc

	return nil
}

// Method PathsChannel gets the channel containing paths to the files that will be processed.
func (p Pipeline) PathsChannel() chan string {
	return p.pathsChannel
}

// Method setPathsChannel sets the channel containing paths to the files that will be processed.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setPathsChannel(pathsChannel chan string) error {
	p.pathsChannel = pathsChannel
	return nil
}

// Method CommandChannel gets the the channel that will signal to start the pipeline.
func (p Pipeline) CommandChannel() chan bool {
	return p.commandChannel
}

// Method setCommandChannel sets the channel that will signal to start the pipeline.
// If there is an error, an error is returned, otherwise nil.
func (p *Pipeline) setCommandChannel(commandChannel chan bool) error {
	p.commandChannel = commandChannel
	return nil
}

// Method Start initiates processing in the pipeline.
func (p Pipeline) Start() error {

	var wg sync.WaitGroup

	// Recursively scan the path for files to process...
	go p.loadPathsToChannel(p.Path(), p.PathsChannel(), p.CommandChannel(), &wg)

	// Start the loading of paths into the paths channel...
	p.CommandChannel() <- true

	// Spin off a goroutine to process each file in the channel
	for path := range p.PathsChannel() {
		go func() {
			fmt.Printf("\nProcessing: %s", path)

			// Do something... Pass the something along to the next step.
		}()
	}

	// Wait for all goroutines to complete.
	wg.Wait()

	return nil
}

// Method Abort abends processing in the pipeline.
func (p Pipeline) Abort() error {
	// TODO
	return nil
}

// Method Stop terminates processing after all steps have been completed.
func (p Pipeline) Stop() error {
	// TODO
	return nil
}
