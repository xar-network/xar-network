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

package types

type MemBackend struct {
	items []interface{}

	readCh chan chan interface{}
	pubCh  chan interface{}
	quitCh chan bool
}

func NewMemBackend() *MemBackend {
	m := &MemBackend{
		items:  make([]interface{}, 0),
		readCh: make(chan chan interface{}),
		pubCh:  make(chan interface{}),
		quitCh: make(chan bool),
	}

	return m
}

func (m *MemBackend) Start() {
	go func() {
		var waitingReaders []chan interface{}

		for {
			select {
			case reader := <-m.readCh:
				if len(m.items) == 0 {
					waitingReaders = append(waitingReaders, reader)
					continue
				}

				reader <- m.items[0]
				m.items = m.items[1:]
			case item := <-m.pubCh:
				if len(waitingReaders) > 0 {
					for _, reader := range waitingReaders {
						reader <- item
					}
					waitingReaders = make([]chan interface{}, 0)
					continue
				}

				m.items = append(m.items, item)
			case <-m.quitCh:
				return
			}
		}
	}()
}

func (m *MemBackend) Stop() {
	m.quitCh <- true
}

func (m *MemBackend) Publish(item interface{}) error {
	m.pubCh <- item
	return nil
}

func (m *MemBackend) Consume() interface{} {
	ch := make(chan interface{})
	m.readCh <- ch
	return <-ch
}
