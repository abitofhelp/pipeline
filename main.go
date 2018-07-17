////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package main contains the entry point to the application and configuration settings.
package main

import (
	"fmt"
	. "github.com/abitofhelp/go-helpers/error"
	"github.com/abitofhelp/pipeline/cmd"
	"gopkg.in/urfave/cli.v2"
	"os"
)

const (
	kApplicationFailure = 0
	kApplicationSuccess = 1
)

// Function main is the entry point to the application.
func main() {

	app := &cli.App{
		Name:  "pipeline",
		Usage: "Scans a directory path for image files and injects them into a pipeline for processing.",
		Commands: []*cli.Command{
			&cmd.Start,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enables debug logging",
			},
		},
	}

	err := app.Run(os.Args)
	if IsError(err, func(err error) {
		fmt.Printf("The application encountered an error: %v", err)
		os.Exit(kApplicationFailure)
	}) {
		fmt.Println("The application has completed successfully.")
		os.Exit(kApplicationSuccess)
	}
}
