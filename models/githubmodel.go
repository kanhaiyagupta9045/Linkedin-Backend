package models

import "time"

type Github struct {
	Title         string    `json:"title" validate:"reuired,min=2,max=15"`
	Description   string    `json:"description" validate:"required min=10,max=100"`
	Language_used string    `json:"language" validate:"required"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Project_id    string    `json:"project_id"`
	Readme        string    `json:"readme"`
}
