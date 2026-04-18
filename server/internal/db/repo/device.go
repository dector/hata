package repo

import (
	"context"
	"errors"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/device"
	"hata/internal/orm/house"
	"hata/internal/orm/houserole"
)

var ErrDeviceNotFound = errors.New("device not found")

// DeviceRepo implements the DeviceRepository interface.
type DeviceRepo struct {
	client *orm.Client
}

// Create creates a new device for the given house.
func (r *DeviceRepo) Create(ctx context.Context, houseID, id, name, integrationID string, integrationData *string, state string) (*DeviceData, error) {
	create := r.client.Device.Create().
		SetHouseID(houseID).
		SetDeviceID(id).
		SetName(name).
		SetIntegrationID(integrationID).
		SetState(state)

	if integrationData != nil {
		create = create.SetIntegrationData(*integrationData)
	}

	d, err := create.Save(ctx)
	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("device %q already exists in house %q", id, houseID)
		}
		return nil, fmt.Errorf("failed creating device: %w", err)
	}

	return deviceDataFromEnt(d), nil
}

// ListByHouse lists devices for a specific house.
func (r *DeviceRepo) ListByHouse(ctx context.Context, houseID string) ([]*DeviceData, error) {
	devices, err := r.client.Device.Query().
		Where(device.HouseIDEQ(houseID)).
		Order(device.ByDeviceID()).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed listing devices for house %q: %w", houseID, err)
	}

	results := make([]*DeviceData, len(devices))
	for i, d := range devices {
		results[i] = deviceDataFromEnt(d)
	}

	return results, nil
}

// ListByUser lists devices for all houses the user belongs to.
func (r *DeviceRepo) ListByUser(ctx context.Context, userID int) ([]*DeviceData, error) {
	devices, err := r.client.Device.Query().
		Where(device.HasHouseWith(house.HasRolesWith(houserole.UserIDEQ(userID)))).
		Order(device.ByHouseID(), device.ByDeviceID()).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed listing devices for user %d: %w", userID, err)
	}

	results := make([]*DeviceData, len(devices))
	for i, d := range devices {
		results[i] = deviceDataFromEnt(d)
	}

	return results, nil
}

// GetByHouseAndID retrieves a device by house and device ID.
func (r *DeviceRepo) GetByHouseAndID(ctx context.Context, houseID, id string) (*DeviceData, error) {
	d, err := r.client.Device.Query().
		Where(device.HouseIDEQ(houseID), device.DeviceIDEQ(id)).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed fetching device %q in house %q: %w", id, houseID, err)
	}
	return deviceDataFromEnt(d), nil
}

// UpdateState updates the device state.
func (r *DeviceRepo) UpdateState(ctx context.Context, houseID, id, state string) error {
	count, err := r.client.Device.Update().
		Where(device.HouseIDEQ(houseID), device.DeviceIDEQ(id)).
		SetState(state).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed updating device state: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q in house %q", ErrDeviceNotFound, id, houseID)
	}
	return nil
}

func deviceDataFromEnt(d *orm.Device) *DeviceData {
	return &DeviceData{
		ID:              d.DeviceID,
		Name:            d.Name,
		IntegrationID:   d.IntegrationID,
		IntegrationData: d.IntegrationData,
		State:           d.State,
		HouseID:         d.HouseID,
	}
}
