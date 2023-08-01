package rekuest

import "fmt"

var (
	// ErrInvalidOption happens when user tries to provide an option with
	// invalid value
	ErrInvalidOption = fmt.Errorf("httprequest: invalid parameter")

	// ErrTimeout happens when the API calls timeout
	ErrTimeout = fmt.Errorf("httprequest: timeout")
)
