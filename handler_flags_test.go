package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateFlag(t *testing.T) {
	validFeatureFlag := func() FeatureFlag {
		return FeatureFlag{
			Name:     "test-feature-flag",
			Enabled:  false,
			Rollout:  50,
			Variants: []string{"controlled", "detailed"},
		}
	}

	validFeatureFlagUpdate := func() FeatureFlagUpdate {
		return FeatureFlagUpdate{Name: "test-feature-flag", Enabled: new(true)}
	}

	flagHandler := FlagHandler{store: CreateFlagStore()}
	t.Run("test CreateFlag func", func(t *testing.T) {

		t.Run("Create Flag Normally", func(t *testing.T) {
			// create featureFlag1
			featureFlag1 := validFeatureFlag()
			body1, _ := json.Marshal(featureFlag1)
			req1 := httptest.NewRequest("POST", "/flags", bytes.NewReader(body1))
			appRes1, err1 := flagHandler.CreateFlag(req1)

			if err1 != nil {
				t.Fatalf("expected no error, got error %v", err1.Error())
			}
			if appRes1.Status != http.StatusCreated {
				t.Fatalf("expected status %v, got  %v", http.StatusCreated, appRes1.Status)
			}
			f1 := appRes1.Body.(FeatureFlag)
			if f1.Name != "test-feature-flag" {
				t.Fatalf("expected FeatureFlag name %v, got %v", "test-feature-flag", f1.Name)
			}

			// create featureFlag2
			featureFlag2 := validFeatureFlag()
			featureFlag2.Name = "test-feature-flag-2"
			body2, _ := json.Marshal(featureFlag2)
			req2 := httptest.NewRequest("POST", "/flags", bytes.NewReader(body2))
			appRes2, err2 := flagHandler.CreateFlag(req2)
			if err2 != nil {
				t.Fatalf("expected no error, got error %v", err2.Error())
			}
			if appRes2.Status != http.StatusCreated {
				t.Fatalf("expected status %v, got  %v", http.StatusCreated, appRes2.Status)
			}
			f2 := appRes2.Body.(FeatureFlag)
			if f2.Name != "test-feature-flag-2" {
				t.Fatalf("expected FeatureFlag name %v, got %v", "test-feature-flag-2", f2.Name)
			}
		})

		t.Run("Create Flag Conflict", func(t *testing.T) {
			body, _ := json.Marshal(validFeatureFlag())
			req := httptest.NewRequest("POST", "/flags", bytes.NewReader(body))
			_, err := flagHandler.CreateFlag(req)

			if err == nil {
				t.Fatal("expected error conflict, got no error")
			}

			var appErr *AppError
			errors.As(err, &appErr)
			if appErr.Status != http.StatusConflict {
				t.Fatalf("Expected %v, got %v", http.StatusConflict, appErr.Status)
			}
		})
	})

	t.Run("test UpdateFlag func", func(t *testing.T) {
		t.Run("Update Flag Notfound", func(t *testing.T) {
			body, _ := json.Marshal(FeatureFlagUpdate{Name: "not-exist-feature-flag", Rollout: new(60)})
			req := httptest.NewRequest("POST", "/flags/not-exist-feature-flag", bytes.NewReader(body))
			req.SetPathValue("name", "not-exist-feature-flag")

			_, err := flagHandler.UpdateFlag(req)
			if err == nil {
				t.Fatalf("expected error, got no error")
			}
			var appErr *AppError
			errors.As(err, &appErr)
			if appErr.Status != http.StatusNotFound {
				t.Fatalf("expected status %v, got %v", http.StatusNotFound, appErr.Status)
			}
		})

		t.Run("Update Flag", func(t *testing.T) {
			update := validFeatureFlagUpdate()
			body, _ := json.Marshal(update)
			req := httptest.NewRequest("POST", fmt.Sprintf("/flag/%v", update.Name), bytes.NewReader(body))
			req.SetPathValue("name", update.Name)

			appRes, err := flagHandler.UpdateFlag(req)
			if err != nil {
				t.Fatalf("expected no error, got error %v", appRes)
			}
			if appRes.Status != http.StatusNoContent {
				t.Fatalf("expected status %v, got %v", http.StatusNoContent, appRes.Status)
			}
		})
	})

	t.Run("test ListAllFlag func", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/flags", nil)
		appRes, err := flagHandler.ListAllFlag(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err.Error())
		}
		if appRes.Status != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, appRes.Status)
		}
		flags := appRes.Body.([]FeatureFlag)
		if len(flags) != 2 {
			t.Fatalf("expected flags_len %v, got %v", 2, len(flags))
		}
	})

	t.Run("test Evaluate func", func(t *testing.T) {
		featureFlag1 := validFeatureFlag()
		req := httptest.NewRequest("GET", fmt.Sprintf("/flag/%v/evaluate?user_id=tom", featureFlag1.Name), nil)
		req.SetPathValue("name", featureFlag1.Name)
		appRes, err := flagHandler.Evaluate(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err.Error())
		}
		if appRes.Status != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, appRes.Status)
		}
		answer := appRes.Body.(*FlagEvaluateAnswer)
		if answer.UserID != "tom" {
			t.Fatalf("expected user_id %v, get %v", "tom", answer.UserID)
		}
	})

}
