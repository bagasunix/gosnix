package responses

import "time"

type DeviceGPSResponse struct {
	ID            string     `json:"id"`
	IMEI          string     `json:"imei"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	Protocol      string     `json:"protocol"`
	IsActive      int        `json:"is_active"`
	InstalledAt   time.Time  `json:"installed_at"`
	UninstalledAt *time.Time `json:"uninstalled_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// Untuk response yang lebih detail tentang status device
type DeviceStatusResponse struct {
	DeviceGPSResponse
	VehicleInfo   *VehicleShortInfo `json:"vehicle_info,omitempty"`
	LastLocation  *LocationResponse `json:"last_location,omitempty"`
	OnlineStatus  int               `json:"online_status"`
	LastHeartbeat *time.Time        `json:"last_heartbeat,omitempty"`
}
