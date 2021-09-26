package middleware

import (
	"database/sql"
	"encoding/json" // Package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // Used to access the request and response object of the api
	"strconv"  // Package used to covert string into int type

	"tasklist-api/drivers" // Database connection logic
	"tasklist-api/models"  // Where database schema is defined

	"github.com/gorilla/mux" // Http router
)

// API Calls //

// API Call to get all Tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tasks, err := getAllTasks()

	if err != nil {
		log.Printf("Unable to get all tasks. %v", err)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// API Call to get one Task by ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		return
	}

	task, err := getTask(int64(id))

	if err != nil {
		log.Printf("Unable to get user. %v", err)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// API Call to get one Task by Status
func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		return
	}

	task, err := getAllTasksStatus(int64(id))

	if err != nil {
		log.Printf("Unable to get user. %v", err)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// API Call to Update one task by Status
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		return
	}

	var task models.Tasks

	err = json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		return
	}

	updatedRows := updateTask(int64(id), task)

	msg := fmt.Sprintf("Task updated successfully. Total rows/record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// API Call to Add a new Task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task models.TasksBase

	// decode the json request to task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		return
	}

	insertID := insertTask(task)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Task created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		return
	}

	deletedRows := deleteTask(int64(id))

	msg := fmt.Sprintf("Tasks updated successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DB Calls

// DB Call to get one Task by ID
func getTask(id int64) (models.Tasks, error) {
	db := drivers.CreateConnection()

	defer db.Close()

	var task models.Tasks

	sqlStatement := `select task_id,task_name,task_desc,username,status_name,priority_name,tasks.created_timestamp,tasks.modified_timestamp
					from tasks

					inner join users
						on assigned_to = user_id
					inner join priority
						on priority = priority_id
					inner join status
						on status = status_id
						where task_id=$1;
					`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&task.Task_id, &task.Task_name, &task.Task_desc, &task.Username, &task.Status_name, &task.Priority_name, &task.Created_timestamp, &task.Modified_timestamp)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return task, nil
	case nil:
		return task, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return task, err
}

// DB Call to get all Tasks by Status
func getAllTasksStatus(id int64) ([]models.Tasks, error) {
	db := drivers.CreateConnection()

	defer db.Close()

	var tasks []models.Tasks

	sqlStatement := `select task_id,task_name,task_desc,username,status_name,priority_name,tasks.created_timestamp,tasks.modified_timestamp
					from tasks

					inner join users
						on assigned_to = user_id
					inner join priority
						on priority = priority_id
					inner join status
						on status = status_id
						where status=$1;
					`

	rows, err := db.Query(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Tasks

		err = rows.Scan(&task.Task_id, &task.Task_name, &task.Task_desc, &task.Username, &task.Status_name, &task.Priority_name, &task.Created_timestamp, &task.Modified_timestamp)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		tasks = append(tasks, task)

	}

	return tasks, err
}

// DB Call to get all Tasks
func getAllTasks() ([]models.Tasks, error) {
	db := drivers.CreateConnection()

	defer db.Close()

	var tasks []models.Tasks

	sqlStatement := `select task_id,task_name,task_desc,username,status_name,priority_name,tasks.created_timestamp,tasks.modified_timestamp
					from tasks

					inner join users
						on assigned_to = user_id
					inner join priority
						on priority = priority_id
					inner join status
						on status = status_id;
					`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Tasks

		err = rows.Scan(&task.Task_id, &task.Task_name, &task.Task_desc, &task.Username, &task.Status_name, &task.Priority_name, &task.Created_timestamp, &task.Modified_timestamp)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		tasks = append(tasks, task)

	}

	return tasks, err
}

// DB Call to Update Task
func updateTask(id int64, task models.Tasks) int64 {
	db := drivers.CreateConnection()

	defer db.Close()

	sqlStatement := `
					UPDATE tasks
					SET task_name=$2,
					    task_desc=$3,
					    assigned_to=$4,
					    status=$5,
					    priority=$6,
					    modified_timestamp=now()
					WHERE task_id=$1;
					`

	res, err := db.Exec(sqlStatement, id, task.Task_name, task.Task_desc, task.Username, task.Status_name, task.Priority_name)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// DB Call to add a new task
func insertTask(task models.TasksBase) int64 {
	db := drivers.CreateConnection()

	defer db.Close()

	sqlStatement := `
					insert into tasks (task_name, task_desc, priority, status, assigned_to, created_timestamp, modified_timestamp)
					values($1, $2, 2, 1, 2, now(), now())
					returning task_id;
					`

	var task_id int64

	err := db.QueryRow(sqlStatement, task.Task_name, task.Task_desc).Scan(&task_id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", task_id)

	return task_id
}

// DB Call to delete a task
func deleteTask(id int64) int64 {
	db := drivers.CreateConnection()

	defer db.Close()

	sqlStatement := `delete from tasks where task_id=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
