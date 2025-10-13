package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Linkedin string `json:"linkedin"`
	Email    string `json:"email"`
	JobRole  Role   `json:"job_role"`
}

type Role struct {
	ID          int    `json:"id"`
	Designation string `json:"designation"`
	Description string `json:"description"`
}
