package main

import (
	"context"
	"fmt"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/utils"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/gridio/test-assignment/internal"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/backend"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/identity"
)

const (
	userID = "xxxyyyzzz"
)

func main() {
	username, _ := os.LookupEnv("USERNAME")
	password, _ := os.LookupEnv("PASSWORD")

	ctx := context.Background()
	logger := logrus.WithField("origin", "test-connector")
	apiClient := utils.NewAPIClient("https://eapi.charge.space/api/v5/", logger)

	// First log in to charge amps and get access tokens
	id, err := identity.Login(logger, apiClient, ctx, username, password)
	if err != nil {
		logger.WithError(err).Error("failed to create identity")

		os.Exit(1)
	}

	secretAgent := internal.NewSecretAgent(id.String())

	newIntegrationFactory := backend.Factory(logger, apiClient)

	// 1. Now create new integration instance
	newIntegration := newIntegrationFactory(userID, secretAgent)

	// 2. Get all devices
	devList, err := newIntegration.DoDeviceListRequest(ctx)
	if err != nil {
		logger.WithError(err).Error("error getting device list")

		os.Exit(1)
	}

	logger.WithField("device_list", devList).Info("device list")

	// 3. Get the status of the first device
	if len(devList) < 1 {
		logger.Info("no devices found")

		os.Exit(0)
	}

	devID := devList[0].ID

	charger, ok := newIntegration.(internal.ChargerBackend)
	if !ok {
		panic("charger backend not implemented")
	}

	status, err := charger.DoChargerStatusRequest(ctx, devID)
	if err != nil {
		logger.WithError(err).Error("error getting device status")

		os.Exit(1)
	}

	logger.WithField("status", fmt.Sprintf("%+v", status)).Info("got device status")
}
