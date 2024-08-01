package handlers

import (
	"FinUslugi/internal/models"
	"FinUslugi/internal/services"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func CreateMaterial(db *pgxpool.Pool, logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) {
	var material models.Material
	if err := json.NewDecoder(r.Body).Decode(&material); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	service := services.NewMaterialService(db, logger)
	newMaterial, err := service.CreateMaterial(material.TypeMaterials, material.Status, material.Title, material.Content)
	if err != nil {
		http.Error(w, "Error creating material", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMaterial)
}

func GetMaterialByID(db *pgxpool.Pool, logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request, id string) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	service := services.NewMaterialService(db, logger)
	material, err := service.GetMaterialByID(uuid)
	if err != nil {
		http.Error(w, "Material not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(material)
}

func UpdateMaterial(db *pgxpool.Pool, logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request, id string) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var material models.Material
	if err := json.NewDecoder(r.Body).Decode(&material); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	service := services.NewMaterialService(db, logger)
	updatedMaterial, err := service.UpdateMaterial(uuid, material.Status, material.Title, material.Content)
	if err != nil {
		http.Error(w, "Error update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMaterial)
}

func GetAllMaterials(db *pgxpool.Pool, logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	var typeMaterials *string
	if val, ok := queryParams["type_materials"]; ok && len(val) > 0 {
		typeMaterials = &val[0]
	}

	var createdAtFrom *time.Time
	if val, ok := queryParams["created_at_from"]; ok && len(val) > 0 {
		parsedTime, err := time.Parse(time.RFC3339, val[0])
		if err == nil {
			createdAtFrom = &parsedTime
		}
	}

	var createdAtTo *time.Time
	if val, ok := queryParams["created_at_to"]; ok && len(val) > 0 {
		parsedTime, err := time.Parse(time.RFC3339, val[0])
		if err == nil {
			createdAtTo = &parsedTime
		}
	}

	limit := 10
	if val, ok := queryParams["limit"]; ok && len(val) > 0 {
		parsedLimit, err := strconv.Atoi(val[0])
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			logger.Warnf("Invalid limit format: %v", err)
		}
	}

	offset := 1
	if val, ok := queryParams["offset"]; ok && len(val) > 0 {
		parsedOffset, err := strconv.Atoi(val[0])
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		} else {
			logger.Warnf("Invalid offset format: %v", err)
		}
	}

	service := services.NewMaterialService(db, logger)
	materials, err := service.GetAllMaterials(typeMaterials, createdAtFrom, createdAtTo, limit, offset)
	if err != nil {
		logger.Errorf("Failed to get materials: %v", err)
		http.Error(w, "Error getting materials", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(materials)
}
