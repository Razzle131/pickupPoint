package dto

type ReceptionInfoDto struct {
	Reception ReceptionDto
	Products  []ProductDto
}

type PvzInfoDto struct {
	Pvz        PvzDto
	Receptions []ReceptionInfoDto
}
