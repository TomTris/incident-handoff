package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (incHandler *IncidentHandler) CreateFlag(w http.ResponseWriter, r *http.Request) (*AppResponse, error) {
	f := FeatureFlag{}
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return nil, BadRequest(err)
	}
	if err := f.Validate(); err != nil {
		return nil, BadRequest(err)
	}
	if err := incHandler.FlagStore.Create(f); err != nil {
		if errors.Is(err, ErrFlagAlreadyExist) {
			return nil, Conflict(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusCreated, f), nil
}

func (incHandler *IncidentHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) (*AppResponse, error) {
	u := FeatureFlagUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, BadRequest(err)
	}

	u.Name = r.PathValue("name")
	if err := u.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	if err := incHandler.FlagStore.Update(u); err != nil {
		if errors.Is(err, ErrFlagNotfound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusNoContent, nil), nil
}

func (incHandler *IncidentHandler) ListAllFlag(w http.ResponseWriter, r *http.Request) (*AppResponse, error) {
	allFlags, err := incHandler.FlagStore.AllFlags()
	if err != nil {
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusOK, allFlags), nil
}

func (incHandler *IncidentHandler) Evaluate(w http.ResponseWriter, r *http.Request) (*AppResponse, error) {
	flagName := r.PathValue("name")
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		return nil, BadRequest(errors.New("empty user_id"))
	}

	flagAnswer, err := incHandler.FlagStore.Evaluate(flagName, userID)
	if err != nil {
		if errors.Is(err, ErrFlagNotfound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusOK, flagAnswer), nil
}
