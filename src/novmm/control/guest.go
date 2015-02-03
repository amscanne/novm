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

package control

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

func (control *Control) init() {

	buffer := make([]byte, 1, 1)

	// Read our control byte back.
	n, err := control.proxy.Read(buffer)
	if err != nil {
		// Something went horribly wrong.
		control.client_res <- err
		return
	}
	if n != 1 {
		// We got nothing.
		control.client_res <- InternalGuestError
		return
	}
	if buffer[0] != '?' {
		// This ain't right.
		control.client_res <- InternalGuestError
		return
	}

	// Send our control byte to noguest.
	buffer[0] = '!'
	n, err = control.proxy.Write(buffer)
	if err != nil {
		// Something went horribly wrong.
		control.client_res <- err
		return
	}
	if n != 1 {
		// Can't send anything?
		control.client_res <- InternalGuestError
		return
	}

	// Looks like we're good.
	control.client_res <- nil
}

func (control *Control) barrier() {
	control.client_err = <-control.client_res
	control.client_codec = jsonrpc.NewClientCodec(control.proxy)
	control.client = rpc.NewClientWithCodec(control.client_codec)
}

func (control *Control) Ready() (*rpc.Client, error) {
	control.client_once.Do(control.barrier)
	return control.client, control.client_err
}
