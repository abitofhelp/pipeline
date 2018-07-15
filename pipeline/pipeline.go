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
	"strings"
	"sync"
	"time"
)

// The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC.
var kTimeZero = time.Time{}

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

const (
	kPathsChannelSize = 100
)

const (
	kGoroutineCount = 20
)

// Type IPipeline is an interface that requires implementations of its pipeline methods.
type IPipeline interface {

	// Function Start initiates processing in the pipeline.
	Start() error

	// Function Abort abends processing in the pipeline.
	Abort() error

	// Function Stop terminates processing after all steps have been completed.
	Stop() error

	// Function Build creates the pipeline.
	Build(wg *sync.WaitGroup) error
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
}

// Function New is a factory that creates an initialized Pipeline.
// Parameter path is the path to a folder of images that will be processed.
// Returns an initialized pipeline or error.
func New(path string) (*Pipeline, error) {
	pipeline := &Pipeline{}
	if pipeline == nil {
		return nil, errors.New("failed to create an instance of Pipeline")
	}

	err := pipeline.setPath(path)
	if err != nil {
		return nil, err
	}

	err = pipeline.setStartedUtc(kTimeZero)
	if err != nil {
		return nil, err
	}

	err = pipeline.setEndedUtc(kTimeZero)
	if err != nil {
		return nil, err
	}

	err = pipeline.setStatus(Stopped)
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

// Method Start initiates processing in the pipeline.
func (p Pipeline) Start() error {
	// TODO
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

// Method Build creates the pipeline.
func (p Pipeline) Build(wg *sync.WaitGroup) error {
	defer wg.Done()

	// Create the channel that will provide paths to files for processing.
	pathsChannel := make(chan string, kPathsChannelSize)

	// TODO: Block it from executing until pipleline is set to running.
	// Recursively scan the path for files to process...
	go p.loadPathsToChannel(p.Path(), pathsChannel, wg)

	// Spin off a goroutine to process each file in the channel
	for path := range pathsChannel {
		go func() {
			fmt.Printf("\nProcessing: %s", path)

			// Do something... Pass the something along to the next step.
		}()
	}

	return nil
}
