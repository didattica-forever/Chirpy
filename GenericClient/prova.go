package genericclient

//package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// JSONCRUD provides generic CRUD operations for JSON-serializable structs.
// It uses an in-memory map as a simple key-value store.
type JSONCRUD struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

// NewJSONCRUD creates a new instance of the JSONCRUD library.
func NewJSONCRUD() *JSONCRUD {
	return &JSONCRUD{
		data:  make(map[string][]byte),
		mutex: sync.RWMutex{},
	}
}

// Create marshals a struct to JSON and stores it with a given key.
// The `item` parameter must be a struct or a pointer to one.
func (j *JSONCRUD) Create(key string, item interface{}) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, ok := j.data[key]; ok {
		return errors.New("key already exists")
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item to JSON: %w", err)
	}

	j.data[key] = jsonData
	return nil
}

// Read retrieves a JSON record by key and unmarshals it into the provided struct.
// The `target` parameter must be a pointer to a struct.
func (j *JSONCRUD) Read(key string, target interface{}) error {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	jsonData, ok := j.data[key]
	if !ok {
		return errors.New("key not found")
	}

	err := json.Unmarshal(jsonData, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON into target: %w", err)
	}

	return nil
}

// Update retrieves a record, merges the new data from `item`, and overwrites the old record.
// The `item` parameter must be a struct or a pointer to one.
func (j *JSONCRUD) Update(key string, item interface{}) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	jsonData, ok := j.data[key]
	if !ok {
		return errors.New("key not found")
	}

	// Unmarshal the existing JSON into a generic map to allow for merging
	var existingMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &existingMap); err != nil {
		return fmt.Errorf("failed to unmarshal existing JSON: %w", err)
	}

	// Marshal the new item to JSON and unmarshal it into a generic map
	updateBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal update item: %w", err)
	}
	var updateMap map[string]interface{}
	if err := json.Unmarshal(updateBytes, &updateMap); err != nil {
		return fmt.Errorf("failed to unmarshal update item bytes: %w", err)
	}

	// Merge the new data over the existing data
	for k, v := range updateMap {
		existingMap[k] = v
	}

	// Marshal the merged map back to JSON
	mergedData, err := json.Marshal(existingMap)
	if err != nil {
		return fmt.Errorf("failed to marshal merged data: %w", err)
	}

	j.data[key] = mergedData
	return nil
}

// Delete removes a record by its key.
func (j *JSONCRUD) Delete(key string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, ok := j.data[key]; !ok {
		return errors.New("key not found")
	}

	delete(j.data, key)
	return nil
}

// getAll retrieves all data currently in the store.
func (j *JSONCRUD) getAll() ([][]byte, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if len(j.data) == 0 {
		return nil, nil
	}

	items := make([][]byte, 0, len(j.data))
	for _, v := range j.data {
		items = append(items, v)
	}

	return items, nil
}

// --- REST API Handlers ---

// The server's in-memory key-value store.
var db = NewJSONCRUD()

// createHandler handles POST requests to create a new item.
// It expects a JSON body with "key" and "data" fields.
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key  string          `json:"key"`
		Data json.RawMessage `json:"data"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	if err := db.Create(req.Key, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully created item with key: %s", req.Key)
}

// readHandler handles GET requests to retrieve an item by key.
func readHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[len("/api/data/"):]
	if key == "" {
		http.Error(w, "Key is required in the URL path, e.g., /api/data/mykey", http.StatusBadRequest)
		return
	}

	var target map[string]interface{}
	if err := db.Read(key, &target); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(target); err != nil {
		http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
	}
}

// updateHandler handles PUT requests to update an item.
// It expects a JSON body with "key" and "data" fields.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key  string          `json:"key"`
		Data json.RawMessage `json:"data"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	if err := db.Update(req.Key, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Successfully updated item with key: %s", req.Key)
}

// deleteHandler handles DELETE requests to delete an item by key.
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[len("/api/data/"):]
	if key == "" {
		http.Error(w, "Key is required in the URL path, e.g., /api/data/mykey", http.StatusBadRequest)
		return
	}

	if err := db.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Successfully deleted item with key: %s", key)
}

// getAllHandler handles GET requests to retrieve all items.
func getAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := db.getAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if items == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "[]")
		return
	}

	// Join all JSON objects into a single JSON array
	var result []json.RawMessage
	for _, item := range items {
		result = append(result, json.RawMessage(item))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
	}
}

func main() {
	// Create a new router
	mux := http.NewServeMux()

	// Register handlers for our REST API endpoints
	mux.HandleFunc("/api/data", createHandler) // POST for creation
	mux.HandleFunc("/api/data/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			readHandler(w, r)
		case "PUT":
			updateHandler(w, r)
		case "DELETE":
			deleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/all", getAllHandler) // GET all items

	// Start the server
	log.Println("Starting server on :8080...")
	log.Println("Available endpoints:")
	log.Println(" - POST /api/data to create an item")
	log.Println(" - GET /api/data/{key} to read an item")
	log.Println(" - PUT /api/data/{key} to update an item")
	log.Println(" - DELETE /api/data/{key} to delete an item")
	log.Println(" - GET /api/all to get all items")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}
