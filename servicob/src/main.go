package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/antoniofmoliveira/go-expert-fullcycle-lab1/src/internal/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cep := r.URL.Query().Get("cep")
		temps, status, message, err := usecase.GetWeather(ctx, cep)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, message, status)
			return
		}

		j, err := json.Marshal(temps)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(j))
	})
	http.ListenAndServe(":8081", nil)
}
