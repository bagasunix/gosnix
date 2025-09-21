package responses

import (
	"time"

	"github.com/gofrs/uuid"
)

type DeviceGPSResponse struct {
	ID            uuid.UUID  `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	IMEI          string     `json:"imei" example:"123456789012345"`
	Brand         string     `json:"brand" example:"Garmin"`
	Model         string     `json:"model" example:"GLO 2"`
	Protocol      string     `json:"protocol" example:"NMEA"`
	IsActive      int        `json:"is_active" example:"1"`
	InstalledAt   time.Time  `json:"installed_at" example:"2023-10-01T10:00:00Z"`
	UninstalledAt *time.Time `json:"uninstalled_at,omitempty" example:"2023-12-01T10:00:00Z"`
	CreatedAt     time.Time  `json:"created_at" example:"2023-10-01T10:00:00Z"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" example:"2023-11-01T10:00:00Z"`
}

// Untuk response yang lebih detail tentang status device
type DeviceStatusResponse struct {
	DeviceGPSResponse
	VehicleInfo   *VehicleShortInfo `json:"vehicle_info,omitempty"`
	LastLocation  *LocationResponse `json:"last_location,omitempty"`
	OnlineStatus  int               `json:"online_status"`
	LastHeartbeat *time.Time        `json:"last_heartbeat,omitempty"`
}
