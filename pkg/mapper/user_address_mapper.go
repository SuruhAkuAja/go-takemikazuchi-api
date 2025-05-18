package mapper

import (
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/user_address/dto"
	"strconv"
)

func MapLocationToUserAddress(location dto.UserLocation, userID uint64) model.UserAddress {
	lat, _ := strconv.ParseFloat(location.Lat, 64)
	lon, _ := strconv.ParseFloat(location.Lon, 64)

	return model.UserAddress{
		PlaceId:          strconv.Itoa(location.PlaceID),
		UserId:           userID,
		FormattedAddress: location.DisplayName,
		Route:            location.UserAddress.Industrial, // Bisa disesuaikan dengan data yang sesuai
		Village:          location.UserAddress.Suburb,
		District:         location.UserAddress.Regency,
		City:             location.UserAddress.City,
		Province:         location.UserAddress.State,
		Country:          location.UserAddress.Country,
		PostalCode:       location.UserAddress.Postcode,
		Latitude:         lat,
		Longitude:        lon,
	}
}
