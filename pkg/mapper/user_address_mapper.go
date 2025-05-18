package mapper

import (
	"fmt"
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

func MapUserAddressModelIntoUserAddressDto(userAddressModel *model.UserAddress) *dto.UserAddressResponse {
	var userAddressDto dto.UserAddressResponse
	fmt.Println(userAddressModel)
	userAddressDto.ID = userAddressModel.ID
	userAddressDto.PlaceId = userAddressModel.PlaceId
	userAddressDto.UserId = userAddressModel.UserId
	fmt.Println("1")
	userAddressDto.FormattedAddress = userAddressModel.FormattedAddress
	userAddressDto.AdditionalInformation = userAddressModel.AdditionalInformation
	userAddressDto.StreetNumber = userAddressModel.StreetNumber
	fmt.Println("2")
	userAddressDto.Route = userAddressModel.Route
	userAddressDto.Village = userAddressModel.Village
	userAddressDto.District = userAddressModel.District
	fmt.Println("3")
	userAddressDto.City = userAddressModel.City
	userAddressDto.Province = userAddressModel.Province
	userAddressDto.Country = userAddressModel.Country
	fmt.Println("4")
	userAddressDto.PostalCode = userAddressModel.PostalCode
	userAddressDto.Latitude = userAddressModel.Latitude
	userAddressDto.Longitude = userAddressModel.Longitude
	return &userAddressDto
}
