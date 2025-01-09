package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_services "github.com/alaurentinoofficial/chartai/internal/core/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

func UrlVars(request *http.Request) map[string]string {
	return mux.Vars(request)
}

func HttpHandlerAuthenticated[T any, U any](
	tokenService core_services.TokenService,
	next core_handlers.HandlerAuthenticatedFunc[T, U],
	parseBody bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !tokenService.Validate(authorization) {
			Error(w, core_validations.ErrForbidden)
			return
		}

		claims, err := tokenService.GetClaims(authorization)
		if err != nil {
			Error(w, core_validations.ErrForbidden)
			return
		}

		accountId, ok := claims["accountId"]
		accountUuid, err := uuid.Parse(accountId)
		if !ok || err != nil {
			Error(w, core_validations.ErrForbidden)
			return
		}

		var request *T = new(T)
		if parseBody {
			request, err = ParseBody[T](r.Body)
			if err != nil {
				slog.Error("Error parsing the body")
				Error(w, core_validations.ErrBadRequest)
				return
			}
		}

		urlVars := mux.Vars(r)
		if err := mapstructure.Decode(urlVars, request); err != nil {
			slog.Error(err.Error())
		}

		auth := core_handlers.Authentication{
			AccountId: accountUuid,
		}
		response, err := next(context.Background(), auth, *request)
		if err != nil {
			slog.Error("Error from application: " + err.Error())
			Error(w, err)
			return
		}

		Ok(w, response)
	}
}

func HttpHandler[T any, U any](next core_handlers.HandlerFunc[T, U], parseBody bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)

		var request *T = new(T)
		if parseBody {
			request, err = ParseBody[T](r.Body)
			if err != nil {
				slog.Error("Error parsing the body")
				Error(w, core_validations.ErrBadRequest)
				return
			}
		}

		urlVars := mux.Vars(r)
		if err := mapstructure.Decode(urlVars, request); err != nil {
			slog.Error(err.Error())
		}

		response, err := next(context.Background(), *request)
		if err != nil {
			slog.Error("Error from application: " + err.Error())
			Error(w, err)
			return
		}

		Ok(w, response)
	}
}

func ParseBody[T any](stream io.ReadCloser) (*T, error) {
	var payload T
	err := json.NewDecoder(stream).Decode(&payload)

	return &payload, err
}

func Ok(w http.ResponseWriter, obj any) {
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(obj)
}

type ErrorResponse struct {
	Errors any `json:"errors"`
}

func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/json")

	if formError, ok := err.(*core_validations.FormErrors); ok {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Errors: formError})
		return
	}

	switch err {
	case core_validations.ErrNotFound:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	case core_validations.ErrForbidden:
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	case core_validations.ErrMethodNotAllowed:
		http.Error(w, "Method not allowed", http.StatusForbidden)
		return
	case core_validations.ErrBadRequest:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	slog.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
