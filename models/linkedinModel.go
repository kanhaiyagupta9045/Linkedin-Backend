package models

import "time"

type LinkedInProfile struct {
	FirstName      string            `json:"first_name" validate:"required,min=6,max=50"`
	LastName       string            `json:"last_name" validate:"required,min=6,max=15"`
	Email          string            `json:"email" validate:"email,required"`
	Password       string            `json:"password" validate:"required,min=8,max=50"`
	Phone          string            `json:"phone" validate:"required,min=10,max=10"`
	Location       string            `json:"location" validate:"required"`
	Connections    int               `json:"connections omitempty"`
	Education      []*Education      `json:"education" validate:"required"`
	WorkExperience []*WorkExperience `json:"work_experience,omitempty"`
	Skills         []*string         `json:"skills,omitempty"`
}
type Education struct {
	CollageName  string    `json:"collage_name" validate:"required,min=2,max=100"`
	Degree       string    `json:"degree" validate:"max=50"`
	FieldOfStudy string    `json:"field_of_study" validate:"max=50"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type WorkExperience struct {
	CompanyName string    `json:"company_name" validate:"required,min=2,max=100"`
	Position    string    `json:"position" validate:"required,min=2,max=100"`
	Description string    `json:"description" validate:"required,max=500"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// Path: models/linkedinModel.go
