package api

type GetPvzResponse struct {
	Pvz        PVZ `json:"pvz"`
	Receptions []struct {
		Reception Reception `json:"reception"`
		Products  []Product `json:"products"`
	} `json:"receptions"`
}
