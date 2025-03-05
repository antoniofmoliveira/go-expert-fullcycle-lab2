package usecase

import (
	"context"
	"net/http"
	"testing"
)

func TestGetWeather(t *testing.T) {
	type args struct {
		ctx context.Context
		cep string
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  int
		wantMessage string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "get weather",
			args: args{
				ctx: context.Background(),
				cep: "39408078",
			},
			wantStatus:  http.StatusOK,
			wantMessage: "OK",
			wantErr:     false,
		},
		{
			name: "get invalid zipcode",
			args: args{
				ctx: context.Background(),
				cep: "3940807",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "Unprocessable Entity",
			wantErr:     true,
		},
		{
			name: "get not found zipcode",
			args: args{
				ctx: context.Background(),
				cep: "39408077",
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "Not Found",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotStatus, gotMessage, err := GetWeather(tt.args.ctx, tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetWeather() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("GetWeather() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
