package api

// addStudent adds a new student to Redis Hash
func addStudent(w http.ResponseWriter, r *http.Request) {
	var student Item
	_ = json.NewDecoder(r.Body).Decode(&student)

	ctx := context.Background()

	// Use HSET to set fields in the Hash
	err := redisClient.HSet(ctx, student.ID, "StudentName", student.Data).Err()
	if err != nil {
		http.Error(w, "Error adding student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// getStudentsByName retrieves all students with a specific name from Redis
func getStudentsByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentName := params["name"]

	ctx := context.Background()
	var students []Item

	// Use HSCAN to efficiently iterate over Hash fields with a specific value
	iter := redisClient.HScan(ctx, studentName, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		field := iter.Val()
		value, err := redisClient.HGet(ctx, studentName, field).Result()
		if err != nil {
			http.Error(w, "Error getting student data", http.StatusInternalServerError)
			return
		}

		students = append(students, Item{ID: field, Data: value})
	}

	if err := iter.Err(); err != nil {
		http.Error(w, "Error retrieving students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// endpoint = /students
router.HandleFunc("/students", addStudent).Methods("POST")
// endpoint = /students/name/{name}
router.HandleFunc("/students/name/{name}", getStudentsByName).Methods("GET")
