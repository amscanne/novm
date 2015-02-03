// Copyright 2014 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"noguest/rpc"
	"os"
	"os/exec"
	"syscall"
)

// The default control file.
var control = flag.String("control", "/dev/vport0p0", "control file")

// Should we mount /dev/pts?
var skip_pts = flag.Bool("skip_pts", false, "skip mounting /dev/pts")

func mount(fs string, location string) error {

	// Do we have the location?
	_, err := os.Stat(location)
	if err != nil {
		// Make sure it's a directory.
		err = os.Mkdir(location, 0755)
		if err != nil {
			return err
		}
	}

	// Try to mount it.
	cmd := exec.Command("/bin/mount", "-t", fs, fs, location)
	return cmd.Run()
}

func main() {

	// Parse flags.
	flag.Parse()

	// Open the console.
	console, err := os.OpenFile(*control, os.O_RDWR, 0)
	if err != nil {
		log.Fatal("problem opening console:", err)
	}

	if !*skip_pts {
		// Make sure devpts is mounted.
		err := mount("devpts", "/dev/pts")
		if err != nil {
			log.Fatal("problem mounting /dev/pts:", err)
		}
	}

	// Notify novmm that we're ready.
	buffer := []byte{'?'}
	n, err := console.Write(buffer)
	if err != nil {
		log.Fatal("problem writing to console:", err)
	}
	if n != 1 {
		log.Fatal("nil write to console")
	}

	// Read our response.
	n, err = console.Read(buffer)
	if err != nil {
		log.Fatal("problem reading from console:", err)
	}
	if n != 1 {
		log.Fatal("nil read from console")
	}
	if buffer[0] != '!' {
		log.Fatal("unexpected response:", buffer[0])
	}

	// Since we don't have any init to setup basic
	// things, like our hostname we do some of that here.
	syscall.Sethostname([]byte("novm"))

	// Small victory.
	log.Printf("~~~ NOGUEST ~~~")

	// Create our RPC server.
	rpc.Run(console)
}
