package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"  // package used to covert string into int type

	"tasklist-api/drivers" // drivers package with database connection logic
	"tasklist-api/models"  // models package where User schema is defined

	"github.com/gorilla/mux" // used to get the params from the route
)

// Response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Users //

// API Call to get one User by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

// API Call to get all Users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

// DB Call to get one User by ID
func getUser(id int64) (models.Users, error) {
	// create the postgres db connection
	db := drivers.CreateConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user models.Users

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE user_id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Created_timestamp, &user.Modified_timestamp)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return user, err
}

// DB Call to get all Users.
func getAllUsers() ([]models.Users, error) {
	// create the postgres db connection
	db := drivers.CreateConnection()

	// close the db connection
	defer db.Close()

	var users []models.Users

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.Users

		// unmarshal the row object to user
		err = rows.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Created_timestamp, &user.Modified_timestamp)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}

// Tasks //

// API Call to get all Tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	tasks, err := getAllTasks()

	if err != nil {
		log.Fatalf("Unable to get all tasks. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(tasks)
}

// API Call to get one Task
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	task, err := getTask(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(task)
}

// DB Call to get one Task by ID
func getTask(id int64) (models.Tasks, error) {
	// create the postgres db connection
	db := drivers.CreateConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var task models.Tasks

	// create the select sql query
	sqlStatement := `select task_id,task_name,task_desc,first_name,last_name,status_name,priority_name,tasks.created_timestamp,tasks.modified_timestamp
					from tasks

					inner join users
						on assigned_to = user_id
					inner join priority
						on priority = priority_id
					inner join status
						on status = status_id
						where task_id=$1;
					`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&task.Task_id, &task.Task_name, &task.Task_desc, &task.First_name, &task.Last_name, &task.Status_name, &task.Priority_name, &task.Created_timestamp, &task.Modified_timestamp)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return task, nil
	case nil:
		return task, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return task, err
}

// DB Call to get all Tasks
func getAllTasks() ([]models.Tasks, error) {
	// create the postgres db connection
	db := drivers.CreateConnection()

	// close the db connection
	defer db.Close()

	var tasks []models.Tasks

	// create the select sql query
	sqlStatement := `select task_id,task_name,task_desc,first_name,last_name,status_name,priority_name,tasks.created_timestamp,tasks.modified_timestamp
					from tasks

					inner join users
						on assigned_to = user_id
					inner join priority
						on priority = priority_id
					inner join status
						on status = status_id;
					`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var task models.Tasks

		// unmarshal the row object to user
		err = rows.Scan(&task.Task_id, &task.Task_name, &task.Task_desc, &task.First_name, &task.Last_name, &task.Status_name, &task.Priority_name, &task.Created_timestamp, &task.Modified_timestamp)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		tasks = append(tasks, task)

	}

	// return empty user on error
	return tasks, err
}
