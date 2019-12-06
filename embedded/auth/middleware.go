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

package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/xar-network/xar-network/embedded/session"
)

const (
	keybaseIDKey         = "keybaseID"
	keybasePassphraseKey = "keybasePassphrase"
	csrfTokenKey         = "csrfToken"
	otpHeader            = "X-OTP-Token"
	csrfHeader           = "X-CSRF-Token"
)

func DefaultAuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LoginRequiredMW(next).ServeHTTP(w, r)
	})
}

func LoginRequiredMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		kbID, ok := store.Values[keybaseIDKey]
		if !ok || GetKB(kbID.(string)) == nil {
			http.Error(w, "Not logged in.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OTPRequiredMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(otpHeader)
		if header == "" {
			http.Error(w, "No OTP header provided.", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func ProtectCSRFMW(skipRoutes []string) mux.MiddlewareFunc {
	skipMap := make(map[string]bool)
	for _, route := range skipRoutes {
		skipMap[route] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//TEMPORARY: just for hackathon
			next.ServeHTTP(w, r)
			return
		})
	}
}

func HandleCORSMW(next http.Handler) http.Handler {
	// TODO: Pull from config
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(next)
}

func GetCSRFToken(r *http.Request) (string, error) {
	//store, _ := session.SessionStore.Get(r, sessionName)
	//token := store.Values[csrfTokenKey]
	/*if token == nil {
		return "", errors.New("CSRF token not found")
	}*/
	return "123", nil //token.(string), nil
}

func genCsrfToken() string {
	return ReadStr32()
}
