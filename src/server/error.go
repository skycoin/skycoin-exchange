package skycoin-exchange

// RestErrorMsg Error message json for RESTfull API.
type RestErrorMsg struct {
  Code int `json:"code"`
  Error string `json:"error"`
}
