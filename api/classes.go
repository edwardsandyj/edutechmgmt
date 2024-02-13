// addClass adds a new class to Redis Hash
func addClass(w http.ResponseWriter, r *http.Request) {
	var class Item
	_ = json.NewDecoder(r.Body).Decode(&class)

	ctx := context.Background()

	// Use HSET to set fields in the Hash
	err := redisClient.HSet(ctx, class.ID, "Name", class.Data).Err()
	if err != nil {
		http.Error(w, "Error adding class", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

// getClassesByName retrieves all classes with a specific name from Redis
func getClassesByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	className := params["name"]

	ctx := context.Background()
	var classes []Item

	// Use HSCAN to efficiently iterate over Hash fields with a specific value
	iter := redisClient.HScan(ctx, className, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		field := iter.Val()
		value, err := redisClient.HGet(ctx, className, field).Result()
		if err != nil {
			http.Error(w, "Error getting class data", http.StatusInternalServerError)
			return
		}

		classes = append(classes, Item{ID: field, Data: value})
	}

	if err := iter.Err(); err != nil {
		http.Error(w, "Error retrieving classes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

// getClassesByTeacher retrieves all classes with a specific teacher from Redis
func getClassesByTeacher(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	teacherName := params["teacher"]

	ctx := context.Background()
	var classes []Item

	// Use HSCAN to efficiently iterate over Hash fields with a specific value
	iter := redisClient.HScan(ctx, teacherName, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		field := iter.Val()
		value, err := redisClient.HGet(ctx, teacherName, field).Result()
		if err != nil {
			http.Error(w, "Error getting class data", http.StatusInternalServerError)
			return
		}

		classes = append(classes, Item{ID: field, Data: value})
	}

	if err := iter.Err(); err != nil {
		http.Error(w, "Error retrieving classes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

// getClassesByRoom retrieves all classes in a specific room from Redis
func getClassesByRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomName := params["room"]

	ctx := context.Background()
	var classes []Item

	// Use HSCAN to efficiently iterate over Hash fields with a specific value
	iter := redisClient.HScan(ctx, roomName, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		field := iter.Val()
		value, err := redisClient.HGet(ctx, roomName, field).Result()
		if err != nil {
			http.Error(w, "Error getting class data", http.StatusInternalServerError)
			return
		}

		classes = append(classes, Item{ID: field, Data: value})
	}

	if err := iter.Err(); err != nil {
		http.Error(w, "Error retrieving classes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

// endpoint = /classes
router.HandleFunc("/classes", addClass).Methods("POST")
// endpoint = /classes/name/{name}
router.HandleFunc("/classes/name/{name}", getClassesByName).Methods("GET")
// endpoint = /classes/teacher/{teacher}
router.HandleFunc("/classes/teacher/{teacher}", getClassesByTeacher).Methods("GET")
// endpoint = /classes/room/{room}
router.HandleFunc("/classes/room/{room}", getClassesByRoom).Methods("GET")
