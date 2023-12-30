package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
)

// Fetches data from the Google Workspace Admin API
func getDataFromAdminAPI(ctx context.Context) ([]Item, error) {
	// Configure Google API client
	config, err := google.JWTConfigFromJSON([]byte("YOUR_SERVICE_ACCOUNT_JSON_KEY"), admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryDeviceReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)

	// Fetch users
	users, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	usersResult, err := users.Users.List().Domain("yourdomain.com").Do()
	if err != nil {
		return nil, err
	}

	// Fetch devices
	devices, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	devicesResult, err := devices.Chromeosdevices.List().CustomerId("my_customer").Projection("BASIC").Do()
	if err != nil {
		return nil, err
	}

	// Transform API results into our data model
	var data []Item
	for _, user := range usersResult.Users {
		data = append(data, Item{ID: user.Id, Type: "User", Data: user.PrimaryEmail})
	}

	for _, device := range devicesResult.Chromeosdevices {
		data = append(data, Item{ID: device.DeviceId, Type: "Device", Data: device.Model})
	}

	return data, nil
}
