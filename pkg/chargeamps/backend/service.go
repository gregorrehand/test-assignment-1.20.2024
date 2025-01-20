package backend

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/gridio/test-assignment/internal"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/identity"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/utils"
)

type Identity interface {
	// AccessToken provides access token string to be used for device list and status requests
	AccessToken() string

	// IsUnauthorized should return true if the token is expired or invalid
	IsUnauthorized() bool
}

type Backend struct {
	tokenAgent internal.SecretAgent
	userID     string
	id         Identity
	apiClient  *utils.APIClient
	logger     logrus.FieldLogger

	// Whatever fields you need
}

// 1. Write integration factory function that produces integrations with the following signature

func Factory(log logrus.FieldLogger, apiClient *utils.APIClient) func(string, internal.SecretAgent) internal.DeviceListProvider {
	f := func(userID string, sa internal.SecretAgent) internal.DeviceListProvider {
		bck := Backend{
			tokenAgent: sa,
			userID:     userID,
			id:         identity.CreateFromSecretAgent(log.WithField("user_id", userID), sa),
			logger:     log.WithField("user_id", userID),
			apiClient:  apiClient,
		}

		return &bck
	}

	return f
}

// 2. Implement interface internal.ChargerBackend

func (b *Backend) DoDeviceListRequest(ctx context.Context) ([]internal.DeviceMetadata, error) {
	// TODO implement me
	// Personally, I would not put any identity logic in the backend service, but since the methods were already provided I did not remove them.
	if b.id.IsUnauthorized() {
		return nil, unAuthorizedAccessError("DoDeviceListRequest")
	}

	var devices []internal.DeviceMetadata

	err := b.apiClient.Get(ctx, "chargepoints/owned", b.id.AccessToken(), &devices)
	if err != nil {
		b.logger.Error("Failed to fetch device list: ", err)
		return nil, err
	}

	return devices, nil
}

func (b *Backend) IsUnauthorized() bool {
	// TODO implement me
	return b.id.IsUnauthorized()
}

var errUnauthorized = errors.New("unauthorized access")

func unAuthorizedAccessError(op string) error {
	return fmt.Errorf("UnAuthorizedAccess %w : %s", errUnauthorized, op)
}

func (b *Backend) DoChargerStatusRequest(ctx context.Context, id internal.PhysicalID) (*internal.ChargerStatus, error) {
	// TODO implement me
	if b.id.IsUnauthorized() {
		return nil, unAuthorizedAccessError("DoChargerStatusRequest")
	}

	var status internal.ChargerStatus

	endpoint := "chargepoints/" + string(id) + "/status"
	if err := b.apiClient.Get(ctx, endpoint, b.id.AccessToken(), &status); err != nil {
		b.logger.Error("Failed to fetch charger status: ", err)
		return nil, err
	}

	return &status, nil
}

// Not sure what the idea with these was as remoteStart endpoint also requires rfid logic, which seems to be beyond the scope of this task
func (b *Backend) StartCharge(ctx context.Context, id internal.PhysicalID, p internal.Power) error {
	// TODO implement me
	panic("implement me")

}
func (b *Backend) Stop(ctx context.Context, id internal.PhysicalID) error {
	// TODO implement me
	panic("implement me")

}
