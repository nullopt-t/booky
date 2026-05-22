package inventory

type ReservedResponse struct {
	Reserved int `json:"reserved"`
}

type AvailableResponse struct {
	Available int `json:"available"`
}
