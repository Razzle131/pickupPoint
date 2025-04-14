package consts

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
)

// aliases for api enum types
var (
	ReceptionStatusActive string = string(api.InProgress)
	ReceptionStatusClosed string = string(api.Close)
)

var (
	UserRoleModerator string = string(api.UserRoleModerator)
	UserRoleEmployee  string = string(api.UserRoleEmployee)
)

var (
	PvzCityMoscow       string = string(api.Москва)
	PvzCityKazan        string = string(api.Казань)
	PvzCityStPetersburg string = string(api.СанктПетербург)
)

var (
	ProductTypeShoes   string = string(api.ProductTypeОбувь)
	ProductTypeClothes string = string(api.ProductTypeОдежда)
	ProductTypeDevices string = string(api.ProductTypeЭлектроника)
)

// util constants
const (
	SliceMinCap    int           = 16
	ContextTimeout time.Duration = time.Millisecond * 100
	HttpTimeout    time.Duration = time.Millisecond * 100
)
