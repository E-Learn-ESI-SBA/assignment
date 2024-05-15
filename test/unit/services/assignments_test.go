package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/utils"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Res struct {
	Message string `json:"message"`
}

var globalAssignment models.Assignment

var admin utils.LightUser = utils.LightUser{
	Username: "admin",
	Role:     "Admin",
	Email:    "admin@gmail.com",
	ID:       uuid.New().String(),
}

var teacher1 utils.LightUser = utils.LightUser{
	Username: "mhammed",
	Role:     "Teacher",
	Email:    "f.mhammed@gmail.com",
	ID:       uuid.New().String(),
}

var teacher2 utils.LightUser = utils.LightUser{
	Username: "poysa",
	Role:     "Teacher",
	Email:    "y.poysa@gmail.com",
	ID:       uuid.New().String(),
}

var secretKey string = "A1B2C3D4E5F6G7H8I9J0K"

var err error

var adminToken string
var teacher1Token string
var teacher2Token string

func TestCreateAssignment(t *testing.T) {
	adminToken, err = utils.GenerateToken(admin, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	teacher1Token, err = utils.GenerateToken(teacher1, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	if err != nil {
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	globalAssignment := models.Assignment{
		ID:          uuid.New(),
		Title:       "archi",
		Description: "this is an assignment",
		Deadline:    time.Now(),
		Year:        "1",
		Groups:      []string{uuid.New().String(), uuid.New().String()},
		TeacherId:   teacher1.ID,
		ModuleId:    uuid.New().String(),
	}

	jsonModule, _ := json.Marshal(globalAssignment)
	req, _ := http.NewRequest(
		"POST",
		"http://localhost:8080/assignments/",
		bytes.NewReader(jsonModule),
	)

	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Assignment Created Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	jsonAssignment, _ := json.Marshal(globalAssignment)
	req, _ = http.NewRequest(
		"POST",
		"http://localhost:8080/assignments/",
		bytes.NewReader(jsonAssignment),
	)
	req.Header.Set("Authorization", "Bearer "+teacher1Token)
	//this should succeed
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"message":"Assignment Created Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

}

// func TestGetAssignmentByModuleId(t *testing.T) {
// 	url := "http://localhost:8080/assignments/" + strconv.Itoa(globalAssignment.ID)
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer "+teacher1Token)

// 	res, _ := http.DefaultClient.Do(req)
// 	responseData, _ := io.ReadAll(res.Body)

// 	var resAssignment []models.Assignment
// 	json.Unmarshal(responseData, &resAssignment)
// 	assert.Equal(t, globalAssignment.ID, resAssignment[0].ID)
// }

func TestGetAssignmentById(t *testing.T) {

	url := "http://localhost:8080/assignments/" + globalAssignment.ID.String()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+teacher1Token)

	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)

	var resAssignment models.Assignment
	json.Unmarshal(responseData, &resAssignment)
	assert.Equal(t, globalAssignment.ID, resAssignment.ID)
}

func TestUpdateAssignment(t *testing.T) {

	updatedAssignment := globalAssignment
	updatedAssignment.Title = "updated title..."
	updatedAssignment.Description = "updated description..."

	jsonAssignment, _ := json.Marshal(updatedAssignment)
	req, _ := http.NewRequest(
		"PUT",
		"http://localhost:8080/assignments/"+globalAssignment.ID.String(),
		bytes.NewReader(jsonAssignment),
	)
	req.Header.Set("Authorization", "Bearer "+teacher1Token)

	//this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Assignment Updated Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	// teacher2Token, err = utils.GenerateToken(teacher2, secretKey)
	// if err != nil {
	// 	//throw err and test failed
	// 	log.Printf("Error: %v\n", err)
	// 	panic(err)
	// }

	// //this should return an error
	// req.Header.Set("Authorization", "Bearer "+teacher2Token)
	// res, _ = http.DefaultClient.Do(req)
	// responseData, _ = io.ReadAll(res.Body)
	// mockResponse = `{"error":"Unauthorized"}`
	// assert.Equal(t, mockResponse, string(responseData))

}

func TestDeleteAssignment(t *testing.T) {
	url := "http://localhost:8080/assignments/" + globalAssignment.ID.String()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+teacher2Token)

	// this should succeed
	req.Header.Set("Authorization", "Bearer "+teacher1Token)
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Assignment Deleted Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	// //this should return an error
	// res, _ = http.DefaultClient.Do(req)
	// responseData, _ = io.ReadAll(res.Body)
	// mockResponse = `{"error":"Unauthorized"}`
	// assert.Equal(t, mockResponse, string(responseData))

}
