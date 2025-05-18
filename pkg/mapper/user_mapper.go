package mapper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-viper/mapstructure/v2"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/user/dto"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"golang.org/x/crypto/bcrypt"
	"googlemaps.github.io/maps"
	"net/http"
)

func MapUserDtoIntoUserModel[T *dto.CreateUserDto](userTransferObject T) *model.User {
	var userModel model.User
	err := mapstructure.Decode(userTransferObject, &userModel)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), 14)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
	userModel.Password = string(hashedPassword)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &userModel
}

func MapJwtClaimIntoUserClaim(jwtClaim jwt.MapClaims) (*dto.JwtClaimDto, error) {
	var userClaim dto.JwtClaimDto
	err := mapstructure.Decode(jwtClaim, &userClaim)
	if err != nil {
		return nil, err
	}
	return &userClaim, nil
}

func MapReverseGeocodingIntoUserAddresses(geocodingResult *maps.GeocodingResult,
	userAddress *model.UserAddress,
	userId uint64,
	addressAdditionalAddress string) {
	userAddress.FormattedAddress = geocodingResult.FormattedAddress
	userAddress.PlaceId = geocodingResult.PlaceID
	userAddress.UserId = userId
	userAddress.Longitude = geocodingResult.Geometry.Location.Lng
	userAddress.Latitude = geocodingResult.Geometry.Location.Lat
	userAddress.AdditionalInformation = addressAdditionalAddress
	for _, addressComponent := range geocodingResult.AddressComponents {
		switch addressComponent.Types[0] {
		case "street_address":
		case "street_number":
			userAddress.StreetNumber += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "route":
			userAddress.Route += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "administrative_area_level_4":
			userAddress.Village += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "administrative_area_level_3":
			userAddress.District += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "administrative_area_level_2":
			userAddress.City += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "administrative_area_level_1":
			userAddress.Province += fmt.Sprintf("%s ", addressComponent.LongName)
			break
		case "country":
			userAddress.Country = addressComponent.LongName
			break
		case "postal_code":
			userAddress.PostalCode = addressComponent.LongName
			break
		}
	}
}

func MapUserModelIntoUserDto(userModel *model.User) *dto.UserResponseDto {
	var userDto dto.UserResponseDto
	err := mapstructure.Decode(userModel, &userDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
	userDto.CreatedAt = userModel.CreatedAt.Format("2006-01-02 15:04:05")
	userDto.UpdatedAt = userModel.UpdatedAt.Format("2006-01-02 15:04:05")
	return &userDto
}
