package models

import (
	"database/sql"
	"time"
)

type Device struct {
	ID                int        `json:"id"`
	DeviceName        string     `json:"deviceName"`
	DeviceDescription string     `json:"deviceDescription"`
	Patches           []Patch    `json:"patches"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}

type DeviceModel struct {
	DB *sql.DB
}

const getDeviceByIdQuery string = `
	SELECT *
	FROM devices
	JOIN patches ON patches.device_id = devices.id
	JOIN users ON users.id = patches.user_id
	WHERE devices.id = $1
`

const getDevicesQuery string = `
	SELECT *
	FROM devices
	JOIN patches ON patches.device_id = devices.id
	LIMIT 10
`

func (m *DeviceModel) GetDeviceById(id string) (Device, error) {
	var device Device
	row := m.DB.QueryRow(getDeviceByIdQuery, id)
	err := row.Scan(&device.ID, &device.DeviceName, &device.DeviceDescription)
	return device, err
}

func (m *DeviceModel) GetDevices() ([]Device, error) {
	var devices []Device
	rows, err := m.DB.Query(getDevicesQuery)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&devices)
	return devices, err
}
