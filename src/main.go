/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/piot/flax/src/server"
)

func parseOptions() (int, string, bool) {
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var destinationHost string
	var listenPort int
	var verbose bool
	commandLine.StringVar(&destinationHost, "destination", "127.0.0.1:32000", "The target")
	commandLine.IntVar(&listenPort, "port", 32001, "The port to listen to")
	commandLine.BoolVar(&verbose, "verbose", false, "Verbose output")

	var flagForceColor = commandLine.Bool("color", false, "Enable color output")
	commandLine.Parse(os.Args[1:])
	if *flagForceColor {
		color.NoColor = false
	}
	return listenPort, destinationHost, verbose
}

func main() {
	color.Cyan("flax 0.1 booting\n")
	listenPort, destinationHost, verbose := parseOptions()
	instance, instanceErr := server.New(listenPort, destinationHost, verbose)
	if instanceErr != nil {
		return
	}
	instance.Forever()
}
