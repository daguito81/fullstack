package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{

		{
			inputJSON:    `{"nickname":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			nickname:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"nickname":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Nickname Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"nickname": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Nickname",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to JSON: %v\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
		}

	}

}
