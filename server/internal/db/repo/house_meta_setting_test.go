package repo_test

import (
	"context"
	"errors"
	"testing"

	"hata/internal/db"
	"hata/internal/db/repo"
)

func TestHouseMetaSettingRepo_SetGetListDelete(t *testing.T) {
	ctx := context.Background()
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer dbInst.Close()

	repos := dbInst.Repos()
	if _, err := repos.House().Create(ctx, "H1", "Main Home"); err != nil {
		t.Fatalf("create house: %v", err)
	}

	setting, err := repos.HouseMetaSetting().Set(ctx, "H1", "hata.weather.v1", "enabled", "true")
	if err != nil {
		t.Fatalf("set setting: %v", err)
	}
	if setting.HouseID != "H1" || setting.Scope != "hata.weather.v1" || setting.Key != "enabled" || setting.Value != "true" {
		t.Fatalf("unexpected created setting: %#v", setting)
	}

	setting, err = repos.HouseMetaSetting().Get(ctx, "H1", "hata.weather.v1", "enabled")
	if err != nil {
		t.Fatalf("get setting: %v", err)
	}
	if setting == nil || setting.Value != "true" {
		t.Fatalf("expected enabled=true, got %#v", setting)
	}

	updated, err := repos.HouseMetaSetting().Set(ctx, "H1", "hata.weather.v1", "enabled", "false")
	if err != nil {
		t.Fatalf("update setting: %v", err)
	}
	if updated.ID != setting.ID || updated.Value != "false" {
		t.Fatalf("expected same setting updated to false, got %#v", updated)
	}

	if _, err := repos.HouseMetaSetting().Set(ctx, "H1", "hata.weather.v1", "units", "metric"); err != nil {
		t.Fatalf("set second setting: %v", err)
	}
	settings, err := repos.HouseMetaSetting().ListByScope(ctx, "H1", "hata.weather.v1")
	if err != nil {
		t.Fatalf("list settings: %v", err)
	}
	if len(settings) != 2 || settings[0].Key != "enabled" || settings[1].Key != "units" {
		t.Fatalf("expected settings ordered by key, got %#v", settings)
	}

	if err := repos.HouseMetaSetting().Delete(ctx, "H1", "hata.weather.v1", "enabled"); err != nil {
		t.Fatalf("delete setting: %v", err)
	}
	missing, err := repos.HouseMetaSetting().Get(ctx, "H1", "hata.weather.v1", "enabled")
	if err != nil {
		t.Fatalf("get deleted setting: %v", err)
	}
	if missing != nil {
		t.Fatalf("expected deleted setting to be missing, got %#v", missing)
	}
}

func TestHouseMetaSettingRepo_DeleteMissing(t *testing.T) {
	ctx := context.Background()
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer dbInst.Close()

	repos := dbInst.Repos()
	if _, err := repos.House().Create(ctx, "H1", "Main Home"); err != nil {
		t.Fatalf("create house: %v", err)
	}

	err = repos.HouseMetaSetting().Delete(ctx, "H1", "hata.weather.v1", "enabled")
	if !errors.Is(err, repo.ErrHouseMetaSettingNotFound) {
		t.Fatalf("expected ErrHouseMetaSettingNotFound, got %v", err)
	}
}
