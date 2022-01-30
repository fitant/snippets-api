package snippet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/contract"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Create(svc service.Service, cfg *config.HTTPServerConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.Context().Value(contract.CS).(contract.CreateSnippet)

		snippet, err := svc.CreateSnippet(data.Snippet, data.Metadata.Language, data.Metadata.Ephemeral)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Printf("%+v\n", snippet)

		if responseType == "raw" {
			w.Write([]byte(snippet.Snippet))
			return
		}

		raw, _ := json.Marshal(snippet)
		w.Write(raw)
	}
}
