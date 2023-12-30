// addDevice adds a new device to Redis Hash
func addDevice(w http.ResponseWriter, r *http.Request) {
	var device Item
	_ = json.NewDecoder(r.Body).Decode(&device)

	ctx := context.Background()

	// Use HSET to set fields in the Hash
	err := redisClient.HSet(ctx, device.ID, "SerialNumber", device.Data).Err()
	if err != nil {
		http.Error(w, "Error adding device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// getDevicesBySerialNumber retrieves all devices with a specific serial number from Redis
func getDevicesBySerialNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serialNumber := params["serialnumber"]

	ctx := context.Background()
	var devices []Item

	// Use HSCAN to efficiently iterate over Hash fields with a specific value
	iter := redisClient.HScan(ctx, serialNumber, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		field := iter.Val()
		value, err := redisClient.HGet(ctx, serialNumber, field).Result()
		if err != nil {
			http.Error(w, "Error getting device data", http.StatusInternalServerError)
			return
		}

		devices = append(devices, Item{ID: field, Data: value})
	}

	if err := iter.Err(); err != nil {
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// endpoint = /devices
router.HandleFunc("/devices", addDevice).Methods("POST")
// endpoint = /devices/serialnumber/{serialnumber}
router.HandleFunc("/devices/serialnumber/{serialnumber}", getDevicesBySerialNumber).Methods("GET")