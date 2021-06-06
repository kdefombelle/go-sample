package common

//APIError defines an generic APIError.
//TODO: used in JwtMiddleware and AuthController for some errors, to be assessed if keep or specialise.
type APIError struct {
	Error            string
	ErrorDescription string
}
