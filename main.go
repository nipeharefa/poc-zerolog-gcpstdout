package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/rs/zerolog/log"
)

type (
	ZerologGCPAdapter struct {
		LogModel LogModel
	}
	LogModel struct {
		Message  string            `json:"message"`
		Severity string            `json:"severity"`
		Labels   map[string]string `json:"logging.googleapis.com/labels"`
	}
)

func (gcp *ZerologGCPAdapter) Write(p []byte) (n int, err error) {
	ss := map[string]string{}
	_ = json.Unmarshal(p, &ss)

	logModel := LogModel{}
	logModel.Severity = ss["level"]
	logModel.Message, _ = ss["message"]

	delete(ss, "level")
	delete(ss, "message")

	logModel.Labels = ss

	by, _ := json.Marshal(logModel)
	fmt.Fprintln(os.Stdout, string(by))

	return len(p), nil
}

func main() {
	w := &ZerologGCPAdapter{}
	log.Logger = log.Output(w)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(5 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		log.Info().
			Str("scope", "payment").
			Str("ss", r.RemoteAddr).
			Str("duitku_error_code", "00").
			Str("duitku_error_desc", "Timeout").
			Msg("Hello World")
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "Hello World")
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", 3000), Handler: router}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
	log.Info().Msg("S")
}
