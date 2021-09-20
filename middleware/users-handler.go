package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"tasklist-api/drivers"
	"tasklist-api/models"

	"github.com/gorilla/mux"
)

// Users //

// API Call to get one User by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

// API Call to get all Users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(users)
}

// DB Call to get one User by ID
func getUser(id int64) (models.Users, error) {
	db := drivers.CreateConnection()

	defer db.Close()

	var user models.Users

	sqlStatement := `SELECT * FROM users WHERE user_id=$1`

	row := db.QueryRow(sqlStatement, id)

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

	return user, err
}

func getAllUsers() ([]models.Users, error) {
	db := drivers.CreateConnection()

	defer db.Close()

	var users []models.Users

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.Users

		err = rows.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Created_timestamp, &user.Modified_timestamp)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)

	}

	return users, err
}
