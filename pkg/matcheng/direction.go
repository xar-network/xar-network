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

package matcheng

import (
	"encoding/json"
	"errors"
)

const (
	Bid Direction = iota
	Ask
)

type Direction uint8

func (d Direction) String() string {
	if d == Bid {
		return "BID"
	}

	return "ASK"
}

func (d *Direction) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	if str == "BID" {
		*d = Bid
	} else if str == "ASK" {
		*d = Ask
	} else {
		return errors.New("invalid direction")
	}

	return nil
}

func (d Direction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}
