package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rodrwan/secretly/internal/database"
)

func RegisterRoutes(router *http.ServeMux, db database.Querier) {
	handler := NewHandler(db)
	// Get all available environments
	router.HandleFunc("GET /api/v1/env", handler.Call(getEnvironments))
	// Create a new environment
	router.HandleFunc("POST /api/v1/env", handler.Call(createEnvironment))
	// Get a specific environment
	router.HandleFunc("GET /api/v1/env/{id}", handler.Call(getEnvironment))
	// Update a specific environment
	router.HandleFunc("PUT /api/v1/env/{id}", handler.Call(updateEnvironment))
	// Delete a specific environment
	router.HandleFunc("DELETE /api/v1/env/{id}", handler.Call(deleteEnvironment))
	// Delete a specific value
	router.HandleFunc("DELETE /api/v1/env/{id}/value/{key}", handler.Call(deleteValue))
}

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

func getEnvironments(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	// Get data from query params, ie: ?name=development
	var envsFromDB []database.Environment
	name := r.URL.Query().Get("name")
	if name != "" {
		env, err := db.GetEnvironmentByName(r.Context(), name)
		if err != nil {
			return Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get environment",
				Error:   err.Error(),
			}, err
		}
		envsFromDB = []database.Environment{env}
	} else {
		envs, err := db.GetAllEnvironments(r.Context())
		if err != nil {
			return Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get environments",
				Error:   err.Error(),
			}, err
		}
		envsFromDB = envs
	}

	envs := make([]Environment, 0)
	for _, env := range envsFromDB {
		valuesFromDB, err := db.GetValuesByEnvironmentID(r.Context(), env.ID)
		if err != nil {
			return Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get values",
				Error:   err.Error(),
			}, err
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

	return Response{
		Code:    http.StatusOK,
		Message: "Environments retrieved",
		Data:    envs,
	}, nil
}

func createEnvironment(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to create environment",
			Error:   err.Error(),
		}, err
	}

	newEnv, err := db.CreateEnvironment(r.Context(), request.Name)
	if err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create environment",
			Error:   err.Error(),
		}, err
	}

	if len(request.Values) > 0 {
		for _, value := range request.Values {
			_, err := db.CreateValue(r.Context(), database.CreateValueParams{
				EnvironmentID: newEnv.ID,
				Key:           value.Key,
				Value:         value.Value,
			})
			if err != nil {
				return Response{
					Code:    http.StatusInternalServerError,
					Message: "Failed to create value",
					Error:   err.Error(),
				}, err
			}
		}
	}

	return Response{
		Code:    http.StatusCreated,
		Message: "Environment created",
		Data:    newEnv,
	}, nil
}

func getEnvironment(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to get environment",
			Error:   err.Error(),
		}, err
	}

	envFromDB, err := db.GetEnvironment(r.Context(), envID)
	if err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get environment",
			Error:   err.Error(),
		}, err
	}

	valuesFromDB, err := db.GetValuesByEnvironmentID(r.Context(), envFromDB.ID)
	if err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get values",
			Error:   err.Error(),
		}, err
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

	return Response{
		Code:    http.StatusOK,
		Message: "Environment retrieved",
		Data:    env,
	}, nil
}

func updateEnvironment(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to update environment",
			Error:   err.Error(),
		}, err
	}

	var request UpdateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to update environment",
			Error:   err.Error(),
		}, err
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
				return Response{
					Code:    http.StatusInternalServerError,
					Message: "Failed to create value",
					Error:   err.Error(),
				}, err
			}
		} else {
			// Si existe, actualizar el valor
			_, err = db.UpdateValue(r.Context(), database.UpdateValueParams{
				ID:    existingValue.ID,
				Value: value.Value,
			})
			if err != nil {
				return Response{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update value",
					Error:   err.Error(),
				}, err
			}
		}
	}

	return Response{
		Code:    http.StatusOK,
		Message: "Environment updated",
		Data:    nil,
	}, nil
}

func deleteEnvironment(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete environment",
			Error:   err.Error(),
		}, err
	}

	err = db.DeleteEnvironment(r.Context(), envID)
	if err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete environment",
			Error:   err.Error(),
		}, err
	}

	return Response{
		Code:    http.StatusOK,
		Message: "Environment deleted",
		Data:    nil,
	}, nil
}

func deleteValue(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error) {
	envID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete value",
			Error:   err.Error(),
		}, err
	}

	// Search for environment
	if _, err := db.GetEnvironment(r.Context(), envID); err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete value",
			Error:   err.Error(),
		}, err
	}

	keyID, err := strconv.ParseInt(r.PathValue("key"), 10, 64)
	if err != nil {
		return Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete value",
			Error:   err.Error(),
		}, err
	}

	err = db.DeleteValue(r.Context(), keyID)
	if err != nil {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete value",
			Error:   err.Error(),
		}, err
	}

	return Response{
		Code:    http.StatusOK,
		Message: "Value deleted",
		Data:    nil,
	}, nil
}
