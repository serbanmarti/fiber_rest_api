package internal

import (
	"fmt"
	"runtime"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
)

type (
	errorLocation struct {
		File string
		Line int
	}
	customError struct {
		Location *errorLocation
		Message  string
		Original error
	}
	BackendError  customError
	DatabaseError customError
)

const (
	ErrBEEmail           = "Error occurred while sending the invite email"
	ErrBEHashSalt        = "Error occurred while generating salt for hashing the given password"
	ErrBEInvalidInvite   = "The given invite token is no longer available"
	ErrBEJwtBad          = "Missing or malformed JWT"
	ErrBEJwtInvalid      = "Invalid or expired JWT"
	ErrBEInvalidPassword = "The given email or password do not match any available data"
	ErrBEMongoIDCast     = "Error occurred while casting MongoDB ObjectID"
	ErrBEMongoIDEmpty    = "The given MongoDB ObjectID is empty"
	ErrBENotActive       = "The account trying to log in is not activated"
	ErrBENotAdmin        = "This logged in account does not have permissions in this section"
	ErrBETimeConversion  = "Error occurred while converting a time field"
	ErrBEUserExists      = "A user account is already associated to this email"

	ErrBEQPInvalidChartType    = "The current request has an invalid or empty chartType query parameter"
	ErrBEQPInvalidDateTime     = "The current request has an invalid or empty date query parameter"
	ErrBEQPInvalidIsInside     = "The current request has an invalid or empty isInside query parameter"
	ErrBEQPInvalidIntervalType = "The current request has an invalid or empty intervalType query parameter"
	ErrBEQPInvalidLocation     = "The current request has an invalid or empty location query parameter"
	ErrBEQPInvalidMobile       = "The current request has an invalid or empty mobile query parameter"
	ErrBEQPInvalidTimezone     = "The current request has an invalid or empty timezone query parameter"
	ErrBEQPMissing             = "The current request is missing one or more query parameters"
	ErrBEQPNoRawOnGate         = "The current request is trying to retrieve non-existing raw data on gates"

	ErrDBCursorClose   = "Error occurred while closing the MongoDB cursor"
	ErrDBCursorIterate = "Error occurred while iterating over the MongoDB cursor"
	ErrDBDecode        = "Error occurred while decoding MongoDB documents"
	ErrDBDelete        = "Error occurred while deleting MongoDB documents"
	ErrDBInsert        = "Error occurred while inserting MongoDB documents"
	ErrDBQuery         = "Error occurred while querying MongoDB documents"
	ErrDBUpdate        = "Error occurred while updating MongoDB documents"
	ErrDBNoData        = "No data found to be grabbed in MongoDB query"
	ErrDBNoUpdate      = "No data found to be updated in MongoDB query"
)

// Create a new BackendError
func NewBackendError(message string, original error, skip int) *BackendError {
	// Generate the error location
	_, file, line, _ := runtime.Caller(skip)

	return &BackendError{
		Location: &errorLocation{
			File: file,
			Line: line,
		},
		Message:  message,
		Original: original,
	}
}

// Error function for error interface
func (e *BackendError) Error() string {
	return e.Message
}

// Create a new DatabaseError
func NewDatabaseError(message string, original error, skip int) *DatabaseError {
	// Generate the error location
	_, file, line, _ := runtime.Caller(skip)

	return &DatabaseError{
		Location: &errorLocation{
			File: file,
			Line: line,
		},
		Message:  message,
		Original: original,
	}
}

// Error function for error interface
func (e *DatabaseError) Error() string {
	return e.Message
}

// Error handler with a custom response for Echo instance
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Create the default values for response code and message
	code := fiber.StatusInternalServerError
	message := "Unhandled error"

	switch e := err.(type) {
	case *fiber.Error: // Handle a Fiber error
		code = e.Code
		message = e.Message

	case validator.ValidationErrors: // Handle a validator error
		code = fiber.StatusBadRequest
		message = ""

		// Customize the response message
		for _, v := range e {
			message += fmt.Sprintf("field validation for '%s' failed on the '%s' tag;", v.Field(), v.ActualTag())
		}

	case *BackendError: // Handle a backend error
		// Log the error
		log.Errorf(
			"BackendError :: File:%s - Line:%d :: %s -> %v", e.Location.File, e.Location.Line, e.Message, e.Original,
		)

		// Construct the response
		message = e.Message

		switch message {
		case ErrBEInvalidPassword, ErrBEJwtInvalid, ErrBENotAdmin:
			code = fiber.StatusUnauthorized
		case ErrBEEmail, ErrBEHashSalt, ErrBEMongoIDCast, ErrBETimeConversion:
			code = fiber.StatusInternalServerError
		case ErrBEInvalidInvite, ErrBEMongoIDEmpty, ErrBEUserExists, ErrBEJwtBad,
			ErrBEQPInvalidChartType, ErrBEQPInvalidDateTime, ErrBEQPInvalidIsInside, ErrBEQPInvalidIntervalType,
			ErrBEQPInvalidLocation, ErrBEQPInvalidMobile, ErrBEQPInvalidTimezone, ErrBEQPMissing, ErrBEQPNoRawOnGate:
			code = fiber.StatusBadRequest
		}

	case *DatabaseError: // Handle a database error
		// Log the error
		log.Errorf(
			"DatabaseError :: File:%s - Line:%d :: %s -> %v", e.Location.File, e.Location.Line, e.Message, e.Original,
		)

		// Construct the response
		message = e.Message

		switch message {
		case ErrDBNoData, ErrDBNoUpdate:
			code = fiber.StatusNotFound
		case ErrDBCursorClose, ErrDBCursorIterate, ErrDBDecode, ErrDBDelete, ErrDBInsert, ErrDBQuery, ErrDBUpdate:
			code = fiber.StatusInternalServerError
		}
	}

	// Send the error response
	if err := c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
		"data":    nil,
	}); err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return nil
}
