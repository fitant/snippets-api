package snippet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/contract"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Create(svc service.Service, cfg *config.HTTPServerConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.Context().Value(contract.CS).(contract.CreateSnippet)

		snippet, err := svc.CreateSnippet(data.Data.Snippet, data.Metadata.Language, data.Metadata.Ephemeral)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.BaseURL, snippet.ID),
		}

		raw, _ := json.Marshal(resp)
		w.Write(raw)
	}
}

func Get(svc service.Service, lgr *zap.Logger, responseType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetID := chi.URLParam(req, "snippetID")

		snippet, err := svc.FetchSnippet(snippetID)
		if err != nil {
			if err == db.ErrNoDocuments {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if responseType == "raw" {
			w.Write([]byte(snippet.Snippet))
			return
		}

		resp := contract.CreateSnippet{
			Data: contract.Data{
				Snippet: snippet.Snippet,
			},
			Metadata: contract.Metadata{
				Language:  snippet.Language,
				Ephemeral: snippet.Ephemeral,
			},
		}

		raw, _ := json.Marshal(resp)
		w.Write(raw)
	}
}
