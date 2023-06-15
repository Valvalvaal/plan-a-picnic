package main

import (
	"fmt"
	"log"
	"net/http"
	"server/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/picnics/", addPicnic)
		v1.GET("/picnics/:picnic_id", readPicnic)
		v1.GET("/picnics/", readAllPicnics)
		v1.PUT("/picnics/:picnic_id", updatePicnic)
		v1.DELETE("/picnics/:picnic_id", deletePicnic)

		v1.POST("/users/", addUser)
		v1.GET("/users/:user_id", readUser)
		v1.GET("/users/", readAllUsers)
		v1.PUT("users/:user_id", updateUser)

		v1.POST("/picnics/:picnic_id/users/:user_id", addUserToPicnic)
		v1.GET("/picnics/:picnic_id/users", readAllUsersOfPicnic)
		v1.GET("/users/:user_id/picnics", readAllPicnicsOfUser)
		// v1.DELETE("/picnics/:picnic_id/users/:user_id", deletePicnicFromUser)
		// v1.DELETE("/users/:user_id/picnics/:picnic_id", deleteUserFromPicnic)

		v1.POST("/food-items/", addFoodItem)
		v1.GET("/food-items/:item_id", readFoodItem)
		v1.GET("/food-items/", readAllFoodItems)
		v1.PUT("/food-items/:item_id", updateFoodItem)
		//v1.DELETE("/food-items/:item_id", deleteFoodItem)

		v1.POST("/contributions/", addContribution)
		v1.GET("/contributions/:contribution_id", readContribution)
		v1.GET("/contributions/", readAllContributions)
		v1.PUT("/contributions/:contribution_id", updateContribution)
		v1.DELETE("/contributions/:contribution_id", deleteContribution)
		// TODO: Crear pruebas en postman, implementar delete, read all contibutions

	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.

	r.Run()

	// picnic routes GET (main page)
	// router.GET("/", picnicIndex)
	// router.GET("/picnics", picnicIndex)
	// router.GET("/picnics/:picnic_id/edit", picnicEdit)

	// // picnic endpoints
	// router.GET("/picnics/:picnic_id", readPicnic)
	// router.GET("/picnics/", readAllPicnics)
	// router.PUT("/picnics/:picnic_id", updatePicnic)
	// router.DELETE("/picnics/:picnic_id", deletePicnic)

	// // users endpoints
	// router.POST("/users/", createUser)
	// router.GET("/users/:picnic_id", readUser) ???
	// router.GET("/picnics/", readAllUsers)
	// router.PUT("/users/:picnic_id", updateUser)
	// router.DELETE("/users/:picnic_id", deleteUser) ??? my version: v1.DELETE("/users/:user_id/picnics/:picnic_id", deleteUserFromPicnic) ???

	// // Contributions endpoints
	// router.POST("/contributions/", createContribution)
	// router.GET("/contributions/:contribution_id", readContribution)
	// router.GET("/contributions/", readAllContributions)
	// router.PUT("/contributions/:contribution_id", updateContribution)
	// router.DELETE("/contributions/:contribution_id", deleteContribution)

	// // FoodItems endpoints
	// router.POST("/food-items/", createFoodItem)
	// router.GET("/food-items/:item_id", readFoodItem)
	// router.GET("/food-items/", readAllFoodItems)
	// router.PUT("/food-items/:item_id", updateFoodItem)
	// router.DELETE("/food-items/:item_id", deleteFoodItem)

}

func addPicnic(c *gin.Context) {

	var json models.Picnic

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("json: ", json)

	success, err := models.CreatePicnic(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func readPicnic(c *gin.Context) {

	// grab the Id of the record we want to retrieve
	id, err := strconv.Atoi(c.Param("picnic_id"))
	checkErr(err)

	picnic, err := models.GetPicnicById(int(id))

	checkErr(err)
	// if the name is blank we can assume nothing is found
	if picnic.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Picnic of that id not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": picnic})
	}
}

func readAllPicnics(c *gin.Context) {

	picnics, err := models.GetPicnics()

	checkErr(err)
	c.JSON(http.StatusOK, gin.H{"data": picnics})
}

func updatePicnic(c *gin.Context) {

	var json models.Picnic

	// grab the Id of the record we want to retrieve

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("picnic_id"))
	json.ID = id
	fmt.Println("json: ", json)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdatePicnic(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deletePicnic(c *gin.Context) {

	picnicId, err := strconv.Atoi(c.Param("picnic_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeletePicnic(picnicId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func addUser(c *gin.Context) {

	var json models.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("json: ", json)

	success, err := models.CreateUser(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func readUser(c *gin.Context) {

	// grab the Id of the record we want to retrieve
	id, err := strconv.Atoi(c.Param("user_id"))
	checkErr(err)

	user, err := models.GetUserById(int(id))

	checkErr(err)
	// if the name is blank we can assume nothing is found
	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User of that id not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func readAllUsers(c *gin.Context) {

	users, err := models.GetUsers()

	checkErr(err)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func updateUser(c *gin.Context) {

	var json models.User

	// grab the Id of the record we want to retrieve

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_1": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("user_id"))
	json.ID = id
	fmt.Println("json: ", json)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_2": "Invalid ID"})
	}

	success, err := models.UpdateUser(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error_3": err})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addUserToPicnic(c *gin.Context) {

	// Get the picnic ID from the request URL parameter
	picnicID, err := strconv.Atoi(c.Param("picnic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picnic ID"})
		return
	}

	// Get the user ID from the request URL parameter
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call AddUserToPicnic
	success, err := models.AddUserToPicnic(userID, picnicID)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func readAllUsersOfPicnic(c *gin.Context) {

	// Get the picnic ID from the request URL parameter
	picnicID, err := strconv.Atoi(c.Param("picnic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picnic ID"})
		return
	}

	picnic, err := models.GetPicnicById(picnicID)
	checkErr(err)

	if picnic.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Picnic of that id not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": picnic})

	}

	// Call a function to retrieve the users by picnic ID from the database
	users, err := models.GetUsersByPicnic(picnic.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})

	// }
}

func readAllPicnicsOfUser(c *gin.Context) {

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := models.GetUserById(userID)
	checkErr(err)

	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User of that id not found"})
	}

	picnics, err := models.GetPicnicsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": picnics})
}

func addFoodItem(c *gin.Context) {

	var json models.FoodItem

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("json: ", json)

	success, err := models.CreateFoodItem(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func readFoodItem(c *gin.Context) {

	// get id of the food item to retrieve

	id, err := strconv.Atoi(c.Param("item_id"))
	checkErr(err)

	foodItem, err := models.GetFoodItemById(int(id))
	checkErr(err)

	if foodItem.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Food item of that id not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": foodItem})
	}
}

func readAllFoodItems(c *gin.Context) {

	foodItems, err := models.GetFoodItems()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food items"})
		return
	}

	if len(foodItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No food items found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": foodItems})
}

func updateFoodItem(c *gin.Context) {

	var json models.FoodItem

	// grab the Id of the record we want to retrieve

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("item_id"))
	json.ID = id
	fmt.Println("json: ", json)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateFoodItem(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func addContribution(c *gin.Context) {

	var json models.Contribution

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("json: ", json)

	success, err := models.CreateContribution(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func readContribution(c *gin.Context) {

	// grab the Id of the record we want to retrieve
	id, err := strconv.Atoi(c.Param("contribution_id"))
	checkErr(err)

	contribution, err := models.GetContributionsOfUserToPicnic(int(id), int(id))

	checkErr(err)
	// if the name is blank we can assume nothing is found
	if contribution.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contribution of that id not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": contribution})
	}
}

func readAllContributions(c *gin.Context) {

	contributions, err := models.GetContributions()

	checkErr(err)
	c.JSON(http.StatusOK, gin.H{"data": contributions})
}

func updateContribution(c *gin.Context) {

	var json models.Contribution

	// grab the Id of the record we want to retrieve & update

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("contribution_id"))
	json.ID = id
	fmt.Println("json: ", json)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateContribution(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteContribution(c *gin.Context) {

	contributionId, err := strconv.Atoi(c.Param("contribution_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteContribution(contributionId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}
