package devicestatus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hata/internal/api"
	"hata/internal/db"
)

const defaultPollInterval = time.Minute

// StartPoller periodically refreshes physical device reachability and state.
func StartPoller(ctx context.Context, repos db.Repositories, controller api.DeviceController, interval time.Duration) {
	if repos == nil || controller == nil {
		return
	}
	if interval <= 0 {
		interval = defaultPollInterval
	}

	go func() {
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				refreshAll(ctx, repos, controller)
				timer.Reset(interval)
			}
		}
	}()
}

func refreshAll(ctx context.Context, repos db.Repositories, controller api.DeviceController) {
	devices, err := repos.Device().ListAll(ctx)
	if err != nil {
		fmt.Printf("Error listing devices for status refresh: %v\n", err)
		return
	}

	for _, device := range devices {
		if ctx.Err() != nil {
			return
		}

		statusCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		status, err := controller.GetStatus(statusCtx, device)
		cancel()

		if err != nil {
			if api.IsDeviceNoAck(err) {
				if updateErr := repos.Device().UpdateAvailability(ctx, device.HouseID, device.ID, "offline"); updateErr != nil && !errors.Is(updateErr, context.Canceled) {
					fmt.Printf("Error marking device %q in house %q offline: %v\n", device.ID, device.HouseID, updateErr)
				}
			}
			continue
		}

		if status.Availability == "" {
			status.Availability = "unknown"
		}
		if status.State == "" {
			if updateErr := repos.Device().UpdateAvailability(ctx, device.HouseID, device.ID, status.Availability); updateErr != nil && !errors.Is(updateErr, context.Canceled) {
				fmt.Printf("Error updating availability for device %q in house %q: %v\n", device.ID, device.HouseID, updateErr)
			}
			continue
		}
		if updateErr := repos.Device().UpdateStatus(ctx, device.HouseID, device.ID, status.State, status.Availability); updateErr != nil && !errors.Is(updateErr, context.Canceled) {
			fmt.Printf("Error updating status for device %q in house %q: %v\n", device.ID, device.HouseID, updateErr)
		}
		if status.LightBrightness != nil || status.LightColorPreset != nil {
			if updateErr := repos.Device().UpdateLight(ctx, device.HouseID, device.ID, status.LightBrightness, status.LightColorPreset); updateErr != nil && !errors.Is(updateErr, context.Canceled) {
				fmt.Printf("Error updating light settings for device %q in house %q: %v\n", device.ID, device.HouseID, updateErr)
			}
		}
	}
}
