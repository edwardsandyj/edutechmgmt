package workspace

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/admin/directory/v1"
)

func GetDataFromAdminAPI(ctx context.Context, serviceAccountJSONPath string) ([]Item, error) {
    // Read the service account JSON key from the file
    serviceAccountJSON, err := ioutil.ReadFile(serviceAccountJSONPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read service account JSON file: %v", err)
    }

    // Configure Google API client with the loaded service account JSON key
    config, err := google.JWTConfigFromJSON(serviceAccountJSON, admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryDeviceReadonlyScope)
    if err != nil {
        return nil, fmt.Errorf("failed to configure Google API client: %v", err)
    }

    client := config.Client(ctx)

    // Fetch users
    users, err := admin.New(client)
    if err != nil {
        return nil, fmt.Errorf("failed to create admin service client: %v", err)
    }

    usersResult, err := users.Users.List().Domain("yourdomain.com").Do()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch users: %v", err)
    }

    // Fetch devices
    devices, err := admin.New(client)
    if err != nil {
        return nil, fmt.Errorf("failed to create admin service client: %v", err)
    }

    devicesResult, err := devices.Chromeosdevices.List().CustomerId("my_customer").Projection("BASIC").Do()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch devices: %v", err)
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
