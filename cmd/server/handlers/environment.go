package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rodrwan/secretly/internal/database"
)

type Environment struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Values []Value `json:"values"`
}

type Value struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Name   string  `json:"name"`
	Values []Value `json:"values"`
}

type UpdateEnvironmentRequest struct {
	EnvironmentID int64   `json:"environment_id"`
	Values        []Value `json:"values"`
}

func RegisterRoutes(router *http.ServeMux, db database.Querier) {
	// Get all available environments
	router.HandleFunc("GET /api/v1/env", getEnvironments(db))
	// Create a new environment
	router.HandleFunc("POST /api/v1/env", createEnvironment(db))
	// Get a specific environment
	router.HandleFunc("GET /api/v1/env/{id}", getEnvironment(db))
	// Update a specific environment
	router.HandleFunc("PUT /api/v1/env/{id}", updateEnvironment(db))
	// Delete a specific environment
	router.HandleFunc("DELETE /api/v1/env/{id}", deleteEnvironment(db))
	// Delete a specific value
	router.HandleFunc("DELETE /api/v1/env/{id}/value/{key}", deleteValue(db))
}

func getEnvironments(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get data from query params, ie: ?name=development
		var envsFromDB []database.Environment
		name := r.URL.Query().Get("name")
		if name != "" {
			env, err := db.GetEnvironmentByName(r.Context(), name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			envsFromDB = []database.Environment{env}
		} else {
			envs, err := db.GetAllEnvironments(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			envsFromDB = envs
		}

		envs := make([]Environment, 0)
		for _, env := range envsFromDB {
			valuesFromDB, err := db.GetValuesByEnvironmentID(r.Context(), env.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			values := make([]Value, 0)
			for _, value := range valuesFromDB {
				values = append(values, Value{
					ID:    value.ID,
					Key:   value.Key,
					Value: value.Value,
				})
			}

			envs = append(envs, Environment{
				ID:     env.ID,
				Name:   env.Name,
				Values: values,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(envs)
	}
}

func createEnvironment(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newEnv, err := db.CreateEnvironment(r.Context(), request.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(request.Values) > 0 {
			for _, value := range request.Values {
				_, err := db.CreateValue(r.Context(), database.CreateValueParams{
					EnvironmentID: newEnv.ID,
					Key:           value.Key,
					Value:         value.Value,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newEnv)
	}
}

func getEnvironment(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		envFromDB, err := db.GetEnvironment(r.Context(), envID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		valuesFromDB, err := db.GetValuesByEnvironmentID(r.Context(), envFromDB.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		values := make([]Value, 0)
		for _, value := range valuesFromDB {
			values = append(values, Value{
				ID:    value.ID,
				Key:   value.Key,
				Value: value.Value,
			})
		}

		env := Environment{
			ID:     envFromDB.ID,
			Name:   envFromDB.Name,
			Values: values,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(env)
	}
}

func updateEnvironment(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var request UpdateEnvironmentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, value := range request.Values {
			// Verificar si el valor ya existe para esta clave en el environment
			existingValue, err := db.GetValueByKey(r.Context(), database.GetValueByKeyParams{
				EnvironmentID: envID,
				Key:           value.Key,
			})

			if err != nil {
				// Si no existe, crear un nuevo valor
				_, err = db.CreateValue(r.Context(), database.CreateValueParams{
					EnvironmentID: envID,
					Key:           value.Key,
					Value:         value.Value,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				// Si existe, actualizar el valor
				_, err = db.UpdateValue(r.Context(), database.UpdateValueParams{
					ID:    existingValue.ID,
					Value: value.Value,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func deleteEnvironment(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.DeleteEnvironment(r.Context(), envID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteValue(db database.Querier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Search for environment
		if _, err := db.GetEnvironment(r.Context(), envID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		keyID, err := strconv.ParseInt(r.PathValue("key"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.DeleteValue(r.Context(), keyID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
