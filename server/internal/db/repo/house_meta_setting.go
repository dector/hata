package repo

import (
	"context"
	"errors"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/housemetasetting"
)

// ErrHouseMetaSettingNotFound is returned when a house metasetting does not exist.
var ErrHouseMetaSettingNotFound = errors.New("house metasetting not found")

// HouseMetaSettingRepo implements the HouseMetaSettingRepository interface.
type HouseMetaSettingRepo struct {
	client *orm.Client
}

// Get retrieves a house metasetting by scope and key.
// Returns nil, nil if the setting does not exist.
func (r *HouseMetaSettingRepo) Get(ctx context.Context, houseID, scope, key string) (*HouseMetaSettingData, error) {
	setting, err := r.client.HouseMetaSetting.Query().
		Where(
			housemetasetting.HouseIDEQ(houseID),
			housemetasetting.ScopeEQ(scope),
			housemetasetting.KeyEQ(key),
		).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed querying house metasetting %q/%q for house %q: %w", scope, key, houseID, err)
	}
	return houseMetaSettingDataFromEnt(setting), nil
}

// ListByScope lists all metasettings for a house and scope.
func (r *HouseMetaSettingRepo) ListByScope(ctx context.Context, houseID, scope string) ([]*HouseMetaSettingData, error) {
	settings, err := r.client.HouseMetaSetting.Query().
		Where(
			housemetasetting.HouseIDEQ(houseID),
			housemetasetting.ScopeEQ(scope),
		).
		Order(orm.Asc(housemetasetting.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing house metasettings for scope %q in house %q: %w", scope, houseID, err)
	}

	results := make([]*HouseMetaSettingData, 0, len(settings))
	for _, setting := range settings {
		results = append(results, houseMetaSettingDataFromEnt(setting))
	}
	return results, nil
}

// Set creates or updates a house metasetting.
func (r *HouseMetaSettingRepo) Set(ctx context.Context, houseID, scope, key, value string) (*HouseMetaSettingData, error) {
	setting, err := r.Get(ctx, houseID, scope, key)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		created, err := r.client.HouseMetaSetting.Create().
			SetHouseID(houseID).
			SetScope(scope).
			SetKey(key).
			SetValue(value).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating house metasetting %q/%q for house %q: %w", scope, key, houseID, err)
		}
		return houseMetaSettingDataFromEnt(created), nil
	}

	updated, err := r.client.HouseMetaSetting.UpdateOneID(setting.ID).
		SetValue(value).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating house metasetting %q/%q for house %q: %w", scope, key, houseID, err)
	}
	return houseMetaSettingDataFromEnt(updated), nil
}

// Delete removes a house metasetting.
func (r *HouseMetaSettingRepo) Delete(ctx context.Context, houseID, scope, key string) error {
	count, err := r.client.HouseMetaSetting.Delete().
		Where(
			housemetasetting.HouseIDEQ(houseID),
			housemetasetting.ScopeEQ(scope),
			housemetasetting.KeyEQ(key),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed deleting house metasetting %q/%q for house %q: %w", scope, key, houseID, err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q/%q for house %q", ErrHouseMetaSettingNotFound, scope, key, houseID)
	}
	return nil
}

func houseMetaSettingDataFromEnt(setting *orm.HouseMetaSetting) *HouseMetaSettingData {
	if setting == nil {
		return nil
	}
	return &HouseMetaSettingData{
		ID:      setting.ID,
		HouseID: setting.HouseID,
		Scope:   setting.Scope,
		Key:     setting.Key,
		Value:   setting.Value,
	}
}
