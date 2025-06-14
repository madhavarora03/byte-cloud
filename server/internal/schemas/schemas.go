package schemas

import "time"

type User struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      *string   `json:"password,omitempty"`
	Provider      string    `json:"provider"`
	ProviderID    *string   `json:"provider_id,omitempty"`
	PictureURL    string    `json:"picture_url"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
