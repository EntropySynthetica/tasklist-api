package models

// User schema of the users table
type Users struct {
	User_id            int64  `json:"user_id"`
	Username           string `json:"username"`
	First_name         string `json:"first_name"`
	Last_name          string `json:"last_name"`
	Created_timestamp  string `json:"created_timestamp"`
	Modified_timestamp string `json:"modified_timestamp"`
}

type Tasks struct {
	Task_id            int64  `json:"task_id"`
	Task_name          string `json:"task_name"`
	Task_desc          string `json:"task_desc"`
	Username           string `json:"username"`
	Status_ID          string `json:"status_id"`
	Status_name        string `json:"status_name"`
	Priority_name      string `json:"priority_name"`
	Created_timestamp  string `json:"created_timestamp"`
	Modified_timestamp string `json:"modified_timestamp"`
}

type TasksBase struct {
	Task_id            int64  `json:"task_id"`
	Task_name          string `json:"task_name"`
	Task_desc          string `json:"task_desc"`
	Priority           string `json:"priority"`
	Status             string `json:"status"`
	Assigned_to        string `json:"assigned_to"`
	Created_timestamp  string `json:"created_timestamp"`
	Modified_timestamp string `json:"modified_timestamp"`
}
