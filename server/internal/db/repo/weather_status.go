package repo

import (
	"context"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/weatherstatus"
)

// WeatherStatusRepo implements the WeatherStatusRepository interface.
type WeatherStatusRepo struct {
	client *orm.Client
}

// GetByHouse retrieves cached weather status for a house.
func (r *WeatherStatusRepo) GetByHouse(ctx context.Context, houseID string) (*WeatherStatusData, error) {
	status, err := r.client.WeatherStatus.Query().
		Where(weatherstatus.HouseIDEQ(houseID)).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed querying weather status for house %q: %w", houseID, err)
	}
	return weatherStatusDataFromEnt(status), nil
}

// UpsertByHouse creates or updates cached weather status for a house.
func (r *WeatherStatusRepo) UpsertByHouse(ctx context.Context, data WeatherStatusData) (*WeatherStatusData, error) {
	existing, err := r.GetByHouse(ctx, data.HouseID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		created, err := r.client.WeatherStatus.Create().
			SetHouseID(data.HouseID).
			SetTemperature(data.Temperature).
			SetTemperatureUnit(data.TemperatureUnit).
			SetConditionCode(data.ConditionCode).
			SetConditionText(data.ConditionText).
			SetConditionIcon(data.ConditionIcon).
			SetNillableHumidityPercent(data.HumidityPercent).
			SetNillableWindSpeed(data.WindSpeed).
			SetNillableWindSpeedUnit(data.WindSpeedUnit).
			SetObservedAt(data.ObservedAt).
			SetUpdatedAt(data.UpdatedAt).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating weather status for house %q: %w", data.HouseID, err)
		}
		return weatherStatusDataFromEnt(created), nil
	}

	update := r.client.WeatherStatus.UpdateOneID(existing.ID).
		SetTemperature(data.Temperature).
		SetTemperatureUnit(data.TemperatureUnit).
		SetConditionCode(data.ConditionCode).
		SetConditionText(data.ConditionText).
		SetConditionIcon(data.ConditionIcon).
		SetObservedAt(data.ObservedAt).
		SetUpdatedAt(data.UpdatedAt)
	if data.HumidityPercent == nil {
		update.ClearHumidityPercent()
	} else {
		update.SetHumidityPercent(*data.HumidityPercent)
	}
	if data.WindSpeed == nil {
		update.ClearWindSpeed()
	} else {
		update.SetWindSpeed(*data.WindSpeed)
	}
	if data.WindSpeedUnit == nil {
		update.ClearWindSpeedUnit()
	} else {
		update.SetWindSpeedUnit(*data.WindSpeedUnit)
	}

	updated, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating weather status for house %q: %w", data.HouseID, err)
	}
	return weatherStatusDataFromEnt(updated), nil
}

func weatherStatusDataFromEnt(status *orm.WeatherStatus) *WeatherStatusData {
	if status == nil {
		return nil
	}
	return &WeatherStatusData{
		ID:              status.ID,
		HouseID:         status.HouseID,
		Temperature:     status.Temperature,
		TemperatureUnit: status.TemperatureUnit,
		ConditionCode:   status.ConditionCode,
		ConditionText:   status.ConditionText,
		ConditionIcon:   status.ConditionIcon,
		HumidityPercent: status.HumidityPercent,
		WindSpeed:       status.WindSpeed,
		WindSpeedUnit:   status.WindSpeedUnit,
		ObservedAt:      status.ObservedAt,
		UpdatedAt:       status.UpdatedAt,
	}
}
