package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)
}

func TestGetHash(t *testing.T) {
	expected := "\xee&\xb0\xddJ\xf7\xe7I\xaa\x1a\x8e\xe3\xc1\n\xe9\x92?a\x89\x80w.G?\x88\x19\xa5Ԕ\x0e\r\xb2z\xc1\x85\xf8\xa0\xe1\xd5\xf8O\x88\xbc\x88\u007f\xd6{\x1472\xc3\x04\xcc_\xa9\xad\x8eoW\xf5\x00(\xa8\xff"
	assert.Equal(t, expected, GetHash("test"))
}
func TestIndex(t *testing.T) {
	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.GET("/", Index)

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected to get status 200")

	// Check the response body is what we expect.
	expected := "Welcome!"
	assert.Contains(t, expected, "Welcome")
}

func Test_counter_with_wrong_method_should_return_status_500(t *testing.T) {
	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.GET("/counter/metrics", HandlerCounter)

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	req, err := http.NewRequest(http.MethodGet, "/counter/metrics", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected to get status 500")
}
func Test_counter_with_right_method_and_empty_input_should_return_status_500(t *testing.T) {
	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.POST("/counter/metrics", HandlerCounter)

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	req, err := http.NewRequest(http.MethodPost, "/counter/metrics", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected to get status 500")
}

func Test_counter_with_right_input_should_return_status_200(t *testing.T) {
	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.POST("/counter/metrics", HandlerCounter)

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	data := `name=test&labels={"labels1":"test1","labels2":"test2"}`

	reader := strings.NewReader(data) //Convert string to reader

	req, err := http.NewRequest("POST", "/counter/metrics", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	assert.Equal(t, http.StatusOK, w.Code, "Expected to get status 200")
}
