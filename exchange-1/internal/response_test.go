package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	errResponse := NewErrorResponse("Test Error Message")

	assert.False(t, errResponse.Status)
	assert.Equal(t, "Test Error Message", errResponse.Message)
}

func TestNewSuccessResponse(t *testing.T) {
	successResponse := NewSuccessResponse("Test Success Message")

	assert.True(t, successResponse.Status)
	assert.Equal(t, "Test Success Message", successResponse.Message)
	assert.Nil(t, successResponse.Data)
}

func TestNewSuccessResponseWithData(t *testing.T) {
	data := struct {
		Name string `json:"name"`
	}{
		Name: "John Doe",
	}

	successResponse := NewSuccessResponse("Test Success Message", WithData(data))

	assert.True(t, successResponse.Status)
	assert.Equal(t, "Test Success Message", successResponse.Message)
	assert.Equal(t, data, successResponse.Data)
}

func TestValidationErrorResponse(t *testing.T) {
	validationError := ValidationError{
		Placement:  "body",
		Detail:     "Invalid input",
		Field:      "username",
		Code:       "1001",
		Expression: "username must be alphanumeric",
		Parameter:  "username",
		TraceId:    "abc123",
	}

	errResponse := ValidationErrorResponse{
		APIResponse: APIResponse{
			Status:  false,
			Message: "Validation Error",
		},
		Errors: []ValidationError{validationError},
	}

	assert.False(t, errResponse.Status)
	assert.Equal(t, "Validation Error", errResponse.Message)
	assert.Equal(t, []ValidationError{validationError}, errResponse.Errors)
}
