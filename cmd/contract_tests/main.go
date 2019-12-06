/*

Copyright 2016 All in Bits, Inc
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

package main

import (
	"fmt"

	"github.com/snikch/goodman/hooks"
	"github.com/snikch/goodman/transaction"
)

func main() {
	// This must be compiled beforehand and given to dredd as parameter, in the meantime the server should be running
	h := hooks.NewHooks()
	server := hooks.NewServer(hooks.NewHooksRunner(h))
	h.BeforeAll(func(t []*transaction.Transaction) {
		fmt.Println("Sleep 5 seconds before all modification")
	})
	h.BeforeEach(func(t *transaction.Transaction) {
		fmt.Println("before each modification")
	})
	h.Before("/version > GET", func(t *transaction.Transaction) {
		fmt.Println("before version TEST")
	})
	h.Before("/node_version > GET", func(t *transaction.Transaction) {
		fmt.Println("before node_version TEST")
	})
	h.BeforeEachValidation(func(t *transaction.Transaction) {
		fmt.Println("before each validation modification")
	})
	h.BeforeValidation("/node_version > GET", func(t *transaction.Transaction) {
		fmt.Println("before validation node_version TEST")
	})
	h.After("/node_version > GET", func(t *transaction.Transaction) {
		fmt.Println("after node_version TEST")
	})
	h.AfterEach(func(t *transaction.Transaction) {
		fmt.Println("after each modification")
	})
	h.AfterAll(func(t []*transaction.Transaction) {
		fmt.Println("after all modification")
	})
	server.Serve()
	defer server.Listener.Close()
	fmt.Print(h)
}
