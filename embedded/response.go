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

package embedded

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func PostProcessResponse(w http.ResponseWriter, cliCtx client.CLIContext, resp interface{}) {
	var result []byte

	switch resp.(type) {
	case []byte:
		result = resp.([]byte)
	default:
		var err error
		if cliCtx.Indent {
			result, err = cliCtx.Codec.MarshalJSONIndent(resp, "", "  ")
		} else {
			result, err = cliCtx.Codec.MarshalJSON(resp)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}
