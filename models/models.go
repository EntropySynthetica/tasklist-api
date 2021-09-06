package models

// User schema of the users table
type Users struct {
	User_id            int64  `json:"user_id"`
	Username           string `json:"name"`
	First_name         string `json:"first_name"`
	Last_name          string `json:"last_name"`
	Created_timestamp  string `json:"created_timestamp"`
	Modified_timestamp string `json:"modified_timestamp"`
}
