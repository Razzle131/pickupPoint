package api

type GetPvzResponse struct {
	Pvz        PVZ
	Receptions []struct {
		Reception Reception
		Products  []Product
	}
}
