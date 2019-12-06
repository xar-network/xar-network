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

package ui

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/xar-network/xar-network/embedded/auth"
)

var tmpl *template.Template
var boxHdlr http.Handler
var uiPaths map[string]bool

func RegisterRoutes(_ context.CLIContext, r *mux.Router, _ *codec.Codec) {
	r.PathPrefix("/").HandlerFunc(uiHandler).Methods("GET")
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl == nil {
		box := packr.NewBox("./build")
		boxHdlr = http.FileServer(box)
		tmplStr, err := box.FindString("/index.html")
		if err != nil {
			panic(err)
		}
		t, err := template.New("entry").Parse(tmplStr)
		if err != nil {
			panic(err)
		}
		tmpl = t
	}

	token, err := auth.GetCSRFToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := uiPaths[r.URL.Path]; ok {
		kb, err := auth.GetKBFromSession(r)
		var uexAddr string
		if err == nil {
			uexAddr = kb.GetAddr().String()
		}

		tmplState := TemplateState{
			CSRFToken:  token,
			UEXAddress: uexAddr,
		}
		err = tmpl.Execute(w, tmplState)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	boxHdlr.ServeHTTP(w, r)
}

type TemplateState struct {
	CSRFToken  string
	UEXAddress string
}

func init() {
	uiPaths = make(map[string]bool)
	uiPaths["/exchange"] = true
	uiPaths["/wallet"] = true
	uiPaths["/"] = true
}
