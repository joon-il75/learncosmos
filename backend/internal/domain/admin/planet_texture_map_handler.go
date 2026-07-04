package admin

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h *AdminHandler) publicActivePlanetTextureMapHandler(c *gin.Context) {
	record, err := h.getActivePlanetTextureMap(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, gin.H{"planet_texture_map": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load active planet texture map"})
		return
	}

	resp, err := buildPlanetTextureMapAssetResponse(*record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read planet texture map metadata"})
		return
	}
	if !resp.Exists {
		c.JSON(http.StatusOK, gin.H{"planet_texture_map": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"planet_texture_map": buildPublicPlanetTextureMapAssetResponse(resp)})
}

func (h *AdminHandler) publicListActivePlanetTextureMapsHandler(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		       cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		       created_at, updated_at
		FROM planet_texture_maps
		WHERE is_active = true
		ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active planet texture maps"})
		return
	}
	defer rows.Close()

	maps := make([]publicPlanetTextureMapAssetResponse, 0)
	for rows.Next() {
		var record planetTextureMapRecord
		if err := rows.Scan(
			&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
			&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
			&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan planet texture map"})
			return
		}

		resp, buildErr := buildPlanetTextureMapAssetResponse(record)
		if buildErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read planet texture map metadata"})
			return
		}
		if resp.Exists {
			maps = append(maps, buildPublicPlanetTextureMapAssetResponse(resp))
		}
	}

	c.JSON(http.StatusOK, gin.H{"planet_texture_maps": maps})
}

func (h *AdminHandler) listPlanetTextureMapsHandler(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		       cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		       created_at, updated_at
		FROM planet_texture_maps
		ORDER BY is_active DESC, created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list planet texture maps"})
		return
	}
	defer rows.Close()

	maps := make([]planetTextureMapAssetResponse, 0)
	for rows.Next() {
		var record planetTextureMapRecord
		if err := rows.Scan(
			&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
			&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
			&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan planet texture map"})
			return
		}

		resp, buildErr := buildPlanetTextureMapAssetResponse(record)
		if buildErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read planet texture map metadata"})
			return
		}
		maps = append(maps, resp)
	}

	c.JSON(http.StatusOK, gin.H{"planet_texture_maps": maps})
}

func (h *AdminHandler) createPlanetTextureMapHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, planetTextureMapMaxBytes+4096)

	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	description := strings.TrimSpace(c.PostForm("description"))
	isActive := parsePlanetTypeActive(c.DefaultPostForm("is_active", "true"))
	rotationDuration := parsePlanetTextureRotationDuration(c.PostForm("rotation_duration_seconds"))
	rotationDirection := normalizePlanetTextureRotationDirection(c.PostForm("rotation_direction"))

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	if fileHeader.Size > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isSupportedPlanetTextureUploadExt(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .png, .jpg, .jpeg, .webp texture maps are supported"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(io.LimitReader(file, planetTextureMapMaxBytes+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}
	if int64(len(payload)) > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}

	webpPayload, err := convertPlanetTextureUploadToWebP(c, payload, ext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert texture map to webp"})
		return
	}
	if int64(len(webpPayload)) > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "converted webp file is too large"})
		return
	}

	width, height, err := readWebPDimensions(webpPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid texture map"})
		return
	}
	if err := validatePlanetTextureMapDimensions(width, height); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	textureMapID := uuid.New()
	fileName := textureMapID.String() + ".webp"
	publicURL := "/textures/planets/maps/" + fileName
	targets := []string{
		filepath.Join(planetTextureMapsSourceDirectory(), fileName),
		filepath.Join(planetTextureMapsStandaloneDirectory(), fileName),
	}
	for _, target := range targets {
		if err := writePlanetTypeAssetFile(target, webpPayload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write planet texture map"})
			return
		}
	}

	var record planetTextureMapRecord
	err = h.db.QueryRow(c.Request.Context(), `
		INSERT INTO planet_texture_maps (
			id, name, description, asset_path, width, height, columns, rows, cell_width, cell_height,
			rotation_duration_seconds, rotation_direction, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		          cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		          created_at, updated_at
	`, textureMapID, name, nullableString(description), publicURL, width, height, planetTextureMapColumns, planetTextureMapRows,
		width/planetTextureMapColumns, height/planetTextureMapRows, rotationDuration, rotationDirection, isActive).Scan(
		&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
		&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
		&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create planet texture map"})
		return
	}

	resp, err := buildPlanetTextureMapAssetResponse(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "planet texture map metadata read failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "planet texture map created", "planet_texture_map": resp})
}

func (h *AdminHandler) updatePlanetTextureMapHandler(c *gin.Context) {
	textureMapID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planet texture map id"})
		return
	}

	var req updatePlanetTextureMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	current, err := h.getPlanetTextureMapByID(c, textureMapID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "planet texture map not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load planet texture map"})
		return
	}

	name := current.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
	}
	description := current.Description
	if req.Description != nil {
		description = strings.TrimSpace(*req.Description)
	}
	isActive := current.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	rotationDuration := current.RotationDurationSeconds
	if req.RotationDurationSeconds != nil {
		rotationDuration = clampPlanetTextureRotationDuration(*req.RotationDurationSeconds)
	}
	rotationDirection := current.RotationDirection
	if req.RotationDirection != nil {
		rotationDirection = normalizePlanetTextureRotationDirection(*req.RotationDirection)
	}

	var record planetTextureMapRecord
	err = h.db.QueryRow(c.Request.Context(), `
		UPDATE planet_texture_maps
		SET name = $2, description = $3, is_active = $4, rotation_duration_seconds = $5,
		    rotation_direction = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		          cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		          created_at, updated_at
	`, textureMapID, name, nullableString(description), isActive, rotationDuration, rotationDirection).Scan(
		&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
		&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
		&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet texture map"})
		return
	}

	resp, err := buildPlanetTextureMapAssetResponse(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read planet texture map metadata"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "planet texture map updated", "planet_texture_map": resp})
}

func (h *AdminHandler) replacePlanetTextureMapFileHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, planetTextureMapMaxBytes+4096)

	textureMapID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planet texture map id"})
		return
	}

	current, err := h.getPlanetTextureMapByID(c, textureMapID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "planet texture map not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load planet texture map"})
		return
	}
	if current.ID.String() == defaultPlanetTextureMapID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "default planet texture map file cannot be replaced"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	if fileHeader.Size > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isSupportedPlanetTextureUploadExt(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .png, .jpg, .jpeg, .webp texture maps are supported"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(io.LimitReader(file, planetTextureMapMaxBytes+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}
	if int64(len(payload)) > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}

	webpPayload, err := convertPlanetTextureUploadToWebP(c, payload, ext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert texture map to webp"})
		return
	}
	if int64(len(webpPayload)) > planetTextureMapMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "converted webp file is too large"})
		return
	}

	width, height, err := readWebPDimensions(webpPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid texture map"})
		return
	}
	if err := validatePlanetTextureMapDimensions(width, height); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := uuid.New().String() + ".webp"
	publicURL := "/textures/planets/maps/" + fileName
	targets := []string{
		filepath.Join(planetTextureMapsSourceDirectory(), fileName),
		filepath.Join(planetTextureMapsStandaloneDirectory(), fileName),
	}
	for _, target := range targets {
		if err := writePlanetTypeAssetFile(target, webpPayload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write planet texture map"})
			return
		}
	}

	var record planetTextureMapRecord
	err = h.db.QueryRow(c.Request.Context(), `
		UPDATE planet_texture_maps
		SET asset_path = $2, width = $3, height = $4, columns = $5, rows = $6,
		    cell_width = $7, cell_height = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		          cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		          created_at, updated_at
	`, textureMapID, publicURL, width, height, planetTextureMapColumns, planetTextureMapRows,
		width/planetTextureMapColumns, height/planetTextureMapRows).Scan(
		&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
		&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
		&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet texture map file"})
		return
	}

	resp, err := buildPlanetTextureMapAssetResponse(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read planet texture map metadata"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "planet texture map file replaced", "planet_texture_map": resp})
}

func (h *AdminHandler) deletePlanetTextureMapHandler(c *gin.Context) {
	textureMapID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planet texture map id"})
		return
	}
	current, err := h.getPlanetTextureMapByID(c, textureMapID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "planet texture map not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load planet texture map"})
		return
	}
	if current.ID.String() == defaultPlanetTextureMapID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "default planet texture map cannot be deleted"})
		return
	}

	_, err = h.db.Exec(c.Request.Context(), `DELETE FROM planet_texture_maps WHERE id = $1`, textureMapID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete planet texture map"})
		return
	}

	fileName := filepath.Base(current.AssetPath)
	targets := []string{
		filepath.Join(planetTextureMapsSourceDirectory(), fileName),
		filepath.Join(planetTextureMapsStandaloneDirectory(), fileName),
	}
	for _, target := range targets {
		if removeErr := os.Remove(target); removeErr != nil && !os.IsNotExist(removeErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove planet texture map"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "planet texture map deleted", "id": textureMapID.String()})
}

func (h *AdminHandler) getActivePlanetTextureMap(c *gin.Context) (*planetTextureMapRecord, error) {
	var record planetTextureMapRecord
	err := h.db.QueryRow(c.Request.Context(), `
		SELECT id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		       cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		       created_at, updated_at
		FROM planet_texture_maps
		WHERE is_active = true
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`).Scan(
		&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
		&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
		&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (h *AdminHandler) getPlanetTextureMapByID(c *gin.Context, textureMapID uuid.UUID) (*planetTextureMapRecord, error) {
	var record planetTextureMapRecord
	err := h.db.QueryRow(c.Request.Context(), `
		SELECT id, name, COALESCE(description, ''), asset_path, width, height, columns, rows,
		       cell_width, cell_height, rotation_duration_seconds, rotation_direction, is_active,
		       created_at, updated_at
		FROM planet_texture_maps
		WHERE id = $1
	`, textureMapID).Scan(
		&record.ID, &record.Name, &record.Description, &record.AssetPath, &record.Width, &record.Height,
		&record.Columns, &record.Rows, &record.CellWidth, &record.CellHeight, &record.RotationDurationSeconds,
		&record.RotationDirection, &record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
