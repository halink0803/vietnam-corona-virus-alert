package crawler

// CurrentSituation response from api
type CurrentSituation struct {
	Location  string `json:"location"`
	Confirmed string `json:"confirmed"`
	Deaths    string `json:"deaths"`
	Recovered string `json:"recovered"`
}
