package bterrors

import (
	"encoding/json"
	"net/http"
)

type BTError struct {
	Message string
	ErrCode int
	HTTPError int `json:"-"`
}

func (bt BTError) Error() string {
	return bt.Message
}

func (bt *BTError) JSON() []byte {
	tx, err := json.Marshal(bt)
	if err != nil {
		return []byte("{\"Message\": \"Something went wrong encoding the error!\", \"ErrCode\": -1 }")
	}

	return tx
}

func (s *BTError) Write(w http.ResponseWriter) {
	w.Header().Del("Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.HTTPError)
	w.Write(s.JSON())
}

// WriteError writes the error to the response. Return false if there is no error
func WriteError(w http.ResponseWriter, err error) bool {
	if err != nil {
		cast, ok := err.(*BTError)
		if ok {
			cast.Write(w)
			return true
		}

		NewUnknown(err).Write(w)
		return true
	}

	return false
}

// New creates a new error
func New(msg string, errCode, err int) *BTError {
	return &BTError{
		Message:   msg,
		ErrCode:   errCode,
		HTTPError: err,
	}
}

// NewUnknown can be used to treat random errors
func NewUnknown(err error) *BTError {
	return &BTError{
		Message:   err.Error(),
		ErrCode:   -1,
		HTTPError: http.StatusInternalServerError,
	}
}

// ErrorCantRefresh is sent when refreshing the token did not succeed
var ErrorCantRefresh *BTError = New("Could not refresh the token!", 399, http.StatusUnauthorized)

// ErrorAPI is sent when something went wrong calling the helix twitch's API
var ErrorAPI *BTError = New("Something went wrong calling API!", 400, http.StatusInternalServerError)

// ErrorNoUser is sent when the user is not found on twitch's server
var ErrorNoUser *BTError = New("No twitch user found", 401, http.StatusNotFound)

// ErrorTokenNotFound is sent when a user tries to log in with a non-existant token
var ErrorTokenNotFound *BTError = New("Token not found", 402, http.StatusNotFound)
