package middleware

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fitant/xbin-api/src/view/http/contract"
)

func WithIngestion() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ch := req.Header.Get("Content-Type")
			data := contract.CreateSnippet{
				Metadata: contract.Metadata{
					Ephemeral: true,
				},
			}

			var source io.Reader
			headerParts := strings.Split(ch, ";")
			if len(headerParts) != 0 && headerParts[0] == "multipart/form-data" {
				lang := req.FormValue("language")
				if lang != "" {
					data.Metadata.Language = lang
				} else {
					data.Metadata.Language = "plaintext"
				}

				eph := req.FormValue("ephemeral")
				if eph == "false" {
					data.Metadata.Ephemeral = false
				}

				// If uploaded file has a content type, use that
				f, h, err := req.FormFile("snippet")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if ct := h.Header.Get("Content-Type"); ct != "" {
					data.Metadata.Language = ct
				}

				defer f.Close()
				source = f
			} else {
				source = req.Body
			}

			raw, _ := ioutil.ReadAll(source)

			if len(raw) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if ch == "application/json" {
				err := json.Unmarshal(raw, &data)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if data.Data.Snippet == "" && data.Metadata.Language == "" {
					data.Data.Snippet = string(raw)
					data.Metadata.Language = "application/json"
				}
			} else {
				data.Data.Snippet = string(raw)
				data.Metadata.Language = "plaintext"
			}

			if data.Data.Snippet == "" || data.Metadata.Language == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), contract.CS, data)))
		})
	}
}
