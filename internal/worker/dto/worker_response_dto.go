package dto

type WorkerResponseDto struct {
	Rating               float32 `json:"rating"`
	Revenue              uint32  `json:"revenue"`
	CompletedJobs        uint32  `json:"completed_jobs"`
	Location             string  `json:"location"`
	Availability         bool    `json:"availability"`
	Verified             bool    `json:"verified"`
	EmergencyPhoneNumber string  `json:"emergency_phone_number"`
	CreatedAt            string  `json:"created_at" mapstructure:"-"`
	UpdatedAt            string  `json:"updated_at" mapstructure:"-"`
}
