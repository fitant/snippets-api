package snippet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/contract"
	"github.com/go-chi/chi/v5"
)

func Create(svc service.Service, cfg *config.HTTPServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.Context().Value(contract.CS).(contract.CreateSnippet)

		snippetID, err := svc.CreateSnippet(data.Data.Snippet, data.Metadata.Language, data.Metadata.Ephemeral)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.GetBaseURL(), snippetID),
		}

		raw, _ := json.Marshal(resp)
		req.Header.Add("Content-Type", "application/json")
		w.Write(raw)
	}
}

func Get(svc service.Service, responseType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetID := chi.URLParam(req, "snippetID")

		snippet, err := svc.FetchSnippet(snippetID)
		if err != nil {
			if err == service.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if responseType == "raw" {
			w.Write(snippet.Snippet)
			return
		}

		resp := contract.CreateSnippet{
			Data: contract.Data{
				Snippet: string(snippet.Snippet),
			},
			Metadata: contract.Metadata{
				Language:  "plaintext",
				Ephemeral: snippet.Ephemeral,
			},
		}

		raw, _ := json.Marshal(resp)
		req.Header.Add("Content-Type", "application/json")
		w.Write(raw)
	}
}
