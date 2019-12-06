/*

Copyright 2019 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package testflags

import (
	"flag"
	"testing"
)

var unitTest = flag.Bool("unit", true, "Run unit tests")
var integrationTest = flag.Bool("integration", true, "Run integration tests")

func IntegrationTest(t *testing.T) {
	if !*integrationTest {
		t.SkipNow()
	}

	t.Parallel()
}

func UnitTest(t *testing.T) {
	if !*unitTest && !testing.Short() {
		t.SkipNow()
	}

	t.Parallel()
}
