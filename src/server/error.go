package server

// ErrorMsg Error message json for RESTfull API.
type ErrorMsg struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
