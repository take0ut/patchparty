package models

import (
	"database/sql"
	"time"
)

type Patch struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DeviceID    int        `json:"deviceID"`
	UserID      int        `json:"userID"`
	Downloads   int        `json:"downloads"`
	URI         string     `json:"URI"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type PatchModel struct {
	DB *sql.DB
}

const GetPatchById string = `
	SELECT *
	FROM patches
	JOIN users ON users.id = patches.user_id
	WHERE patches.id = $1
`

const CreatePatch string = `
	INSERT INTO patches (name, description, device_id, user_id, downloadeds, source_url) 
	VALUES ($1, $2, $3, $4, $5, $6)
`

func (m *PatchModel) GetPatchById(id string) (Patch, error) {
	var patch Patch
	row := m.DB.QueryRow(GetPatchById, id)
	err := row.Scan(&patch.ID, &patch.Name, &patch.Description, &patch.DeviceID, &patch.UserID, &patch.Downloads, &patch.URI)
	return patch, err
}

func (m *PatchModel) CreatePatch(patch Patch) (Patch, error) {
	row := m.DB.QueryRow(CreatePatch, patch.Name, patch.Description, patch.DeviceID, patch.UserID, patch.Downloads, patch.URI)
	err := row.Scan(&patch.ID, &patch.Name, &patch.Description, &patch.DeviceID, &patch.UserID, &patch.Downloads, &patch.URI)
	return patch, err
}
