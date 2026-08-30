package model

import "time"

type Profile struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Title       string `json:"title"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
	LinkedInURL string `json:"linkedin_url"`
	GitHubURL   string `json:"github_url"`
	Summary     string `json:"summary"`
}

type SkillGroup struct {
	Category string   `json:"category"`
	Items    []string `json:"items"`
}

type Experience struct {
	ID        int64      `json:"id"`
	Company   string     `json:"company"`
	Client    string     `json:"client,omitempty"`
	Role      string     `json:"role"`
	Location  string     `json:"location"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Bullets   []string   `json:"bullets"`
}

type Project struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Summary string   `json:"summary"`
	Bullets []string `json:"bullets"`
}

type Education struct {
	ID          int64      `json:"id"`
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Location    string     `json:"location"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type Portfolio struct {
	Profile    Profile      `json:"profile"`
	Skills     []SkillGroup `json:"skills"`
	Experience []Experience `json:"experience"`
	Projects   []Project    `json:"projects"`
	Education  []Education  `json:"education"`
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type ContactMessage struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Subject string    `json:"subject"`
	Body    string    `json:"message"`
	Status  string    `json:"status"`
	Created time.Time `json:"created_at"`
}
