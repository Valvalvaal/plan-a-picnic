package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_TABLE_STATEMENTS_SQL = `

CREATE TABLE IF NOT EXISTS users (
  id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name        VARCHAR UNIQUE NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX  IF NOT EXISTS index_users_on_name ON users (name);

CREATE TABLE IF NOT EXISTS picnics (
  id           INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name         VARCHAR NOT NULL,
  location     VARCHAR NOT NULL,
  date         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_picnics (
  id       INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  user_id  INTEGER,
  picnic_id INTEGER,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (picnic_id) REFERENCES picnics(id)
);

CREATE TABLE IF NOT EXISTS food_items (
  id           INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name         VARCHAR NOT NULL,
  measure      VARCHAR NOT NULL,
  url		   VARCHAR NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contributions (
	id 	         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	user_id      INTEGER,
	picnic_id    INTEGER,
	food_item_id INTEGER,
	quantity   INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
  	FOREIGN KEY (picnic_id) REFERENCES picnics(id),
	FOREIGN KEY (food_item_id) REFERENCES food_items(id)
);

`

// A contribution le paso el id de la persona y del picnic
type Picnic struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Date     string `json:"date"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserPicnic struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	PicnicID int `json:"picnic_id"`
}

type FoodItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Measure string `json:"measure"`
}

type Contribution struct {
	ID         int `json:"id"`
	UserID     int `json:"user_id"`
	PicnicID   int `json:"picnic_id"`
	FoodItemID int `json:"food_item_id"`
	Quantity   int `json:"quantity"`
}

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./models/data.db")
	if err != nil {
		return err
	}

	DB = db
	createTables()
	return nil
}

func createTables() {
	_, err := DB.Exec(CREATE_TABLE_STATEMENTS_SQL)
	if err != nil {
		log.Fatal("Error during table creation SQL statements")
		return
	}
}

func GetPicnicById(id int) (Picnic, error) {

	stmt, err := DB.Prepare("SELECT id, name, location, date from picnics WHERE id = ?")

	if err != nil {
		return Picnic{}, err
	}

	picnic := Picnic{}

	sqlErr := stmt.QueryRow(id).Scan(&picnic.ID, &picnic.Name, &picnic.Location, &picnic.Date)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Picnic{}, nil
		}
		return Picnic{}, sqlErr
	}
	return picnic, nil
}

func GetPicnics() ([]Picnic, error) {

	rows, err := DB.Query("SELECT id, name, location, date from picnics")
	picnics := make([]Picnic, 0)
	if err != nil {
		return picnics, err
	}

	for rows.Next() {
		picnic := Picnic{}
		err = rows.Scan(&picnic.ID, &picnic.Name, &picnic.Location, &picnic.Date)

		if err != nil {
			return make([]Picnic, 0), err
		}

		picnics = append(picnics, picnic)
	}

	err = rows.Err()

	if err != nil {
		return make([]Picnic, 0), err
	}

	return picnics, err
}

func CreatePicnic(newPicnic Picnic) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO picnics (name, location, date) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPicnic.Name, newPicnic.Location, newPicnic.Date)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdatePicnic(updatedPicnic Picnic, idToUpdate int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE picnics SET name = ?, location = ?, date = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updatedPicnic.Name, updatedPicnic.Location, updatedPicnic.Date, idToUpdate)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeletePicnic(picnicId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from picnics where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(picnicId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func CreateUser(newUser User) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO users (name) VALUES (?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newUser.Name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetUserById(id int) (User, error) {

	stmt, err := DB.Prepare("SELECT id, name FROM users WHERE id = ?")

	if err != nil {
		return User{}, err
	}

	user := User{}

	sqlErr := stmt.QueryRow(id).Scan(&user.ID, &user.Name)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}
	return user, nil
}

func GetUsers() ([]User, error) {

	rows, err := DB.Query("SELECT id, name FROM users")
	users := make([]User, 0)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Name)

		if err != nil {
			return make([]User, 0), err
		}

		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		return make([]User, 0), err
	}

	return users, err
}

func UpdateUser(updatedUser User, idToUpdate int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET name = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updatedUser.Name, idToUpdate)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func AddUserToPicnic(userID int, picnicID int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO users_picnics (user_id, picnic_id) VALUES (?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userID, picnicID)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetUsersByPicnic(picnicId int) ([]User, error) {
	// Select the necessary data to create a user obj by picnic id from tables users and picnics
	rows, err := DB.Query("SELECT users.id, users.name FROM users INNER JOIN users_picnics ON users.id = users_picnics.user_id WHERE users_picnics.picnic_id = ?", picnicId)
	users := make([]User, 0)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return make([]User, 0), err
		}
		users = append(users, user)
	}
	err = rows.Err()

	if err != nil {
		return make([]User, 0), err
	}

	return users, err

}

func GetPicnicsByUser(userId int) ([]Picnic, error) {
	rows, err := DB.Query("SELECT picnics.id,  picnics.name, picnics.location, picnics.date FROM picnics INNER JOIN users_picnics ON picnics.id = picnic_id WHERE user_id = ?", userId)
	picnics := make([]Picnic, 0)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		picnic := Picnic{}
		err := rows.Scan(&picnic.ID, &picnic.Name, &picnic.Location, &picnic.Date)

		if err != nil {
			return make([]Picnic, 0), err
		}
		picnics = append(picnics, picnic)
	}
	err = rows.Err()
	if err != nil {
		return make([]Picnic, 0), err
	}

	return picnics, err
}

func CreateFoodItem(newFoodItem FoodItem) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Printf("error 1: %v", err)
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO food_items (name, measure, url) VALUES (?, ?, ?)")

	if err != nil {
		fmt.Printf("error 2: %v", err)
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newFoodItem.Name, newFoodItem.Measure, newFoodItem.Url)

	if err != nil {
		fmt.Printf("error 3: %v", err)
		return false, err
	}

	tx.Commit()
	return true, nil
}

func GetFoodItemById(id int) (FoodItem, error) {

	stmt, err := DB.Prepare("SELECT id, name, measure, url FROM food_items WHERE id = ?")

	if err != nil {
		return FoodItem{}, err
	}

	foodItem := FoodItem{}

	sqlErr := stmt.QueryRow(id).Scan(&foodItem.ID, &foodItem.Name, &foodItem.Measure, &foodItem.Url)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return FoodItem{}, nil
		}
		return FoodItem{}, sqlErr
	}
	return foodItem, nil
}

func GetFoodItems() ([]FoodItem, error) {

	rows, err := DB.Query("SELECT id, name, measure, url FROM food_items")
	foodItems := make([]FoodItem, 0)
	if err != nil {
		fmt.Printf("error 1: %v", err)
		return foodItems, err
	}

	for rows.Next() {
		foodItem := FoodItem{}
		err = rows.Scan(&foodItem.ID, &foodItem.Name, &foodItem.Measure, &foodItem.Url)

		if err != nil {
			fmt.Printf("error 2: %v", err)
			return make([]FoodItem, 0), err
		}

		foodItems = append(foodItems, foodItem)
	}

	err = rows.Err()

	if err != nil {
		fmt.Printf("error 3: %v", err)
		return make([]FoodItem, 0), err
	}

	return foodItems, err
}

func UpdateFoodItem(updatedFoodItem FoodItem, idToUpdate int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE food_items SET name = ?, measure = ?, url = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updatedFoodItem.Name, updatedFoodItem.Measure, updatedFoodItem.Url, idToUpdate)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func CreateContribution(newContribution Contribution) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO contributions (user_id, picnic_id, food_item_id, quantity) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newContribution.UserID, newContribution.PicnicID, newContribution.FoodItemID, newContribution.Quantity)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetContributionsOfUserToPicnic(idUser int, idPicnic int) (Contribution, error) {

	stmt, err := DB.Prepare("SELECT id, user_id, picnic_id, food_item_id, quantity from contributions WHERE user_id = ? AND picnic_id = ?")

	if err != nil {
		return Contribution{}, err
	}

	contribution := Contribution{}

	sqlErr := stmt.QueryRow(idUser, idPicnic).Scan(&contribution.ID, &contribution.UserID, &contribution.PicnicID, &contribution.FoodItemID, &contribution.Quantity)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Contribution{}, nil
		}
		return Contribution{}, sqlErr
	}
	return contribution, nil

}

func GetContributions() ([]Contribution, error) {

	rows, err := DB.Query("SELECT id, user_id, picnic_id, food_item_id, quantity FROM contributions")
	contributions := make([]Contribution, 0)
	if err != nil {
		return contributions, err
	}

	for rows.Next() {
		contribution := Contribution{}
		err = rows.Scan(&contribution.ID, &contribution.UserID, &contribution.PicnicID, &contribution.FoodItemID, &contribution.Quantity)

		if err != nil {
			return make([]Contribution, 0), err
		}

		contributions = append(contributions, contribution)
	}

	err = rows.Err()

	if err != nil {
		return make([]Contribution, 0), err
	}
	return contributions, err
}

func UpdateContribution(updatedContribution Contribution, idToUpdate int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE contributions SET user_id = ?, picnic_id = ?, food_item_id = ?, quantity = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updatedContribution.UserID, updatedContribution.PicnicID, updatedContribution.FoodItemID, updatedContribution.Quantity, idToUpdate)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteContribution(contributionId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from contributions where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(contributionId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
