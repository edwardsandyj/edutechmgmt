Make sure to install Redix client package first!!
github.com/go-redis/redis/v8



main.go –
Entry point to the application
Imports functionality for retrieving Google Admin API data from 'workspace'
Implements functionality for interacting with Redis and handling WebSocket connections
Sets up HTTP endpoints for exporting and importing data, editing items, and searching items

redisops.go –
Contains the updateRedisDataStore function which updates Redis based on fetched data
Called from main.go and relies on the Redis client

routing.go – 
Defines HTTP handlers for managing items in Redis
Imports the Redis client and uses it to perform CRUD operations on items

workspace/adminapi.go –
Contains functions for fetching data from the Google Admin API
Requires a service account JSON key to authenticate with the API

devices.go –
Contains functions for adding devices to Redis and retrieving devices by serial number
Relies on the Redis client for interacting with Redis

students.go –
Contains functions for adding students to Redis and retrieving students by name
Relies on the Redis client for interacting with Redis

classes.go –
Contains functions for adding classes to Redis and retrieving classes by class, teacher, and classroom
Relies on the Redis client for interacting with Redis



