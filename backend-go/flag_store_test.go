package main

import (
	"errors"
	"testing"
)

func seefFlagStore() FlagStore {
	return CreateFlagStore()
}

func TestFlagCreate(t *testing.T) {
	flagStore := CreateFlagStore()
	featureFlag := FeatureFlag{
		Name:     "flag1",
		Enabled:  true,
		Rollout:  50,
		Variants: []string{"control", "variant_a"},
	}

	err := flagStore.Create(featureFlag)
	if err != nil {
		t.Fatalf("expect no err, got err %v", err.Error())
	}
	if len(flagStore.Flags) != 1 {
		t.Errorf("expected 1, get %v", len(flagStore.Flags))
	}

	t.Run("recreate flag", func(t *testing.T) {
		err = flagStore.Create(featureFlag)
		if err == nil {
			t.Errorf("flagStore expected error, get 0 error")
		}
	})
}

func TestFlagUpdate(t *testing.T) {
	flagStore := CreateFlagStore()
	featureFlag := FeatureFlag{
		Name:     "dark_mode",
		Enabled:  true,
		Rollout:  50,
		Variants: []string{"control", "variant_a"},
	}

	flagStore.Create(featureFlag)

	t.Run("update enabled", func(t *testing.T) {
		enabled := false
		err := flagStore.Update(FeatureFlagUpdate{Name: "dark_mode", Enabled: &enabled})
		if err != nil {
			t.Fatal(err)
		}
		if flagStore.Flags["dark_mode"].Enabled != false {
			t.Error("expected disabled")
		}
	})

	t.Run("update rollout", func(t *testing.T) {
		rollout := 75
		err := flagStore.Update(FeatureFlagUpdate{Name: "dark_mode", Rollout: &rollout})
		if err != nil {
			t.Fatal(err)
		}
		if flagStore.Flags["dark_mode"].Rollout != 75 {
			t.Errorf("expected 75, got %d", flagStore.Flags["dark_mode"].Rollout)
		}
	})

	t.Run("update both fields", func(t *testing.T) {
		enabled := true
		rollout := 30
		err := flagStore.Update(FeatureFlagUpdate{Name: "dark_mode", Enabled: &enabled, Rollout: &rollout})
		if err != nil {
			t.Fatal(err)
		}
		flag := flagStore.Flags["dark_mode"]
		if flag.Enabled != true {
			t.Error("expected enabled")
		}
		if flag.Rollout != 30 {
			t.Errorf("expected 30, got %d", flag.Rollout)
		}
	})

	t.Run("not found", func(t *testing.T) {
		enabled := true
		err := flagStore.Update(FeatureFlagUpdate{Name: "nonexistent", Enabled: &enabled})
		if !errors.Is(err, ErrFlagNotfound) {
			t.Errorf("expected ErrFlagNotfound, got %v", err.Error())
		}
	})
}

func TestFlagStore_AllFlags(t *testing.T) {
	flagStore := CreateFlagStore()

	t.Run("empty store", func(t *testing.T) {
		flags, err := flagStore.AllFlags()
		if err != nil {
			t.Errorf("expected no error")
		}
		if len(flags) != 0 {
			t.Errorf("expected 0, got %d", len(flags))
		}
	})

	t.Run("returns all", func(t *testing.T) {
		flagStore.Create(FeatureFlag{Name: "a", Variants: []string{"v1"}})
		flagStore.Create(FeatureFlag{Name: "b", Variants: []string{"v1"}})
		flags, err := flagStore.AllFlags()
		if err != nil {
			t.Errorf("expected no error")
		}
		if len(flags) != 2 {
			t.Errorf("expected 2, got %d", len(flags))
		}
	})
}

func TestFlagStore_Evaluate(t *testing.T) {
	fs := CreateFlagStore()
	fs.Create(FeatureFlag{
		Name:     "dark_mode",
		Enabled:  true,
		Rollout:  100,
		Variants: []string{"control", "variant_a"},
	})

	fs.Create(FeatureFlag{
		Name:     "disabled_flag",
		Enabled:  false,
		Rollout:  100,
		Variants: []string{"control"},
	})
	t.Run("flag not found", func(t *testing.T) {
		_, err := fs.Evaluate("nonexistent", "user1")
		if !errors.Is(err, ErrFlagNotfound) {
			t.Errorf("expected ErrFlagNotfound, got %v", err.Error())
		}
	})

	t.Run("empty userID returns not in rollout", func(t *testing.T) {
		answer, err := fs.Evaluate("dark_mode", "")
		if err != nil {
			t.Fatal(err)
		}
		if answer.InRollout != false {
			t.Error("expected not in rollout for empty userID")
		}
		if answer.Variant != nil {
			t.Error("expected nil variant")
		}
	})

	t.Run("disabled flag returns not in rollout", func(t *testing.T) {
		answer, err := fs.Evaluate("disabled_flag", "user1")
		if err != nil {
			t.Fatal(err)
		}
		if answer.InRollout != false {
			t.Error("expected not in rollout for disabled flag")
		}
	})

	t.Run("100% rollout always in rollout", func(t *testing.T) {
		answer, err := fs.Evaluate("dark_mode", "user1")
		if err != nil {
			t.Fatal(err)
		}
		if answer.InRollout != true {
			t.Error("expected in rollout at 100%")
		}
		if answer.Variant == nil {
			t.Fatal("expected non-nil variant")
		}
		if *answer.Variant != "control" && *answer.Variant != "variant_a" {
			t.Errorf("unexpected variant: %s", *answer.Variant)
		}
	})

	t.Run("0% rollout never in rollout", func(t *testing.T) {
		fs.Create(FeatureFlag{
			Name:     "zero_rollout",
			Enabled:  true,
			Rollout:  0,
			Variants: []string{"control"},
		})
		answer, err := fs.Evaluate("zero_rollout", "user1")
		if err != nil {
			t.Fatal(err)
		}
		if answer.InRollout != false {
			t.Error("expected not in rollout at 0%")
		}
	})

	t.Run("deterministic for same user", func(t *testing.T) {
		a1, _ := fs.Evaluate("dark_mode", "user1")
		a2, _ := fs.Evaluate("dark_mode", "user1")
		if a1.InRollout != a2.InRollout {
			t.Error("same user should get same rollout result")
		}
		if a1.InRollout && *a1.Variant != *a2.Variant {
			t.Error("same user should get same variant")
		}
	})

	t.Run("answer fields populated correctly", func(t *testing.T) {
		answer, _ := fs.Evaluate("dark_mode", "user42")
		if answer.Name != "dark_mode" {
			t.Errorf("expected dark_mode, got %s", answer.Name)
		}
		if answer.UserID != "user42" {
			t.Errorf("expected user42, got %s", answer.UserID)
		}
	})
}
