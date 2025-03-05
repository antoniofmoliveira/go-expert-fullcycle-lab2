package usecase

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	neturl "net/url"
	"os"
	"strings"

	"github.com/antoniofmoliveira/go-expert-fullcycle-lab1/src/internal/dto"
)

// getCep retrieves the cep information from the ViaCEP API.
//
// It constructs a request to the ViaCEP API using the given cep and sends an HTTP GET request.
// If the request is successful and the response contains valid data, it unmarshals the JSON response
// into a Viacep struct and converts it to a Cep struct. It returns the Cep instance, along with
// a status code and a message indicating the success or failure of the operation.
//
// Possible return scenarios include:
// - 200, "OK": if the request is successful and the cep is found.
// - 404, "Not Found": if the cep cannot be found or the response indicates an error.
// - 422, "Unprocessable Entity": if the cep is invalid.
// - 408, "Request Timeout": if the request times out.
// - 500, "Internal Server Error": if there is a failure in processing the request or response.
// - 503, "Service Unavailable": if the ViaCEP service is unavailable.
// It returns an error if any network or processing error occurs.

func getCep(ctx context.Context, cep string) (rcep dto.Cep, status int, message string, error error) {
	url := "http://viacep.com.br/ws/{{cep}}/json/"
	url = strings.Replace(url, "{{cep}}", cep, 1)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return dto.Cep{}, 500, "Internal Server Error", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.Cep{}, 500, "Internal Server Error", err
	}

	switch res.StatusCode {
	case http.StatusOK:

		body, error := io.ReadAll(res.Body)
		if error != nil {
			return dto.Cep{}, 500, "Internal Server Error", errors.New("fail to read the body response: " + error.Error())
		}
		defer res.Body.Close()

		// obs this api doesnt return a 404 error if the cep is not found
		// the return is 200 status with a body with "erro": "true"
		if strings.Contains(string(body), `"erro": "true"`) {
			return dto.Cep{}, 404, "Not Found", errors.New("can not find zipcode")
		}

		cepdto, err := dto.NewViacepFromJson(string(body))
		if err != nil {
			return dto.Cep{}, 500, "Internal Server Error", err
		}
		rcep := dto.Cep{
			Cep:          cepdto.Cep,
			State:        cepdto.Uf,
			City:         cepdto.Localidade,
			Neighborhood: cepdto.Bairro,
			Street:       cepdto.Logradouro,
		}
		slog.Info(rcep.City)
		return rcep, 200, "OK", nil
	case http.StatusRequestTimeout:
		return dto.Cep{}, 408, "Request Timeout", errors.New("time exceeded")

	case http.StatusNotFound:
		return dto.Cep{}, 404, "Not Found", errors.New("can not find zipcode")

	case http.StatusBadRequest:
		return dto.Cep{}, 422, "Unprocessable Entity", errors.New("invalid zipcode")

	case http.StatusInternalServerError:
		return dto.Cep{}, 500, "Internal Server Error", errors.New("internal server error")

	case http.StatusServiceUnavailable:
		return dto.Cep{}, 503, "Service Unavailable", errors.New("service unavailable")

	default:
		return dto.Cep{}, 404, "Not Found", errors.New("can not find zipcode")
	}

}

// GetWeather gets the current weather for a given cep.
//
// It first calls getCep to get the cep information, then makes a request to the weather api
// using the city name from the cep information. It then marshalls the response into a TempResponse
// and returns it, along with the appropriate status and message.
//
// It returns an error if the cep is invalid, the request to the weather api fails, or if the
// response from the weather api is invalid.
//
// It returns 200, "OK" if the request is successful, 408, "Request Timeout" if the request times out,
// 404, "Not Found" if the cep is invalid or the weather api can not find the city,
// 400, "Bad Request" if the request is invalid,
// 422, "Unprocessable Entity" if the city is invalid,
// 500, "Internal Server Error" if there is an internal error,
// 503, "Service Unavailable" if the weather api is unavailable,
// or an error if an unknown error occurs.
func GetWeather(ctx context.Context, cep string) (temps dto.TempResponse, status int, message string, eror error) {
	rcep, status, message, err := getCep(ctx, cep)
	if err != nil {
		return dto.TempResponse{}, status, message, err
	}

	api_key := os.Getenv("API_KEY")

	url := "http://api.weatherapi.com/v1/current.json?key={{key}}&q={{city}}&aqi=no"
	url = strings.Replace(url, "{{key}}", api_key, 1)
	url = strings.Replace(url, "{{city}}", neturl.QueryEscape(rcep.City), 1)

	slog.Info(url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return dto.TempResponse{}, 500, "Internal Server Error", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.TempResponse{}, 500, "Internal Server Error", err
	}

	slog.Info("Status: ", "Status Code: ", res.StatusCode)
	switch res.StatusCode {
	case http.StatusOK:
		body, error := io.ReadAll(res.Body)
		if error != nil {
			return dto.TempResponse{}, 500, "Internal Server Error", errors.New("fail to read the body response: " + error.Error())
		}
		defer res.Body.Close()

		tempdto, err := dto.NewWeatherApiFromJson(string(body))
		if err != nil {
			return dto.TempResponse{}, 500, "Internal Server Error", err
		}
		rtemp := dto.TempResponse{
			City:   tempdto.Location.Name,
			Temp_C: tempdto.Current.TempC,
			Temp_F: tempdto.Current.TempC*1.8 + 32,
			Temp_K: tempdto.Current.TempC + 273.0,
		}
		return rtemp, 200, "OK", nil
	case http.StatusRequestTimeout:
		return dto.TempResponse{}, 408, "Request Timeout", errors.New("time exceeded")

	case http.StatusNotFound:
		return dto.TempResponse{}, 404, "Not Found", errors.New("can not find weather")

	case http.StatusBadRequest:
		return dto.TempResponse{}, 400, "Bad Request", errors.New("bad request")

	case http.StatusUnprocessableEntity:
		return dto.TempResponse{}, 422, "Unprocessable Entity", errors.New("invalid city")

	case http.StatusInternalServerError:
		return dto.TempResponse{}, 500, "Internal Server Error", errors.New("internal server error")

	case http.StatusServiceUnavailable:
		return dto.TempResponse{}, 503, "Service Unavailable", errors.New("service unavailable")

	default:
		return dto.TempResponse{}, 500, "Internal Server Error", errors.New("unknown error")
	}

}
