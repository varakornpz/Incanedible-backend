package utils

import (
	"errors"
	"strings"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/varakornpz/mygorm"
)


func GetUUIDByContext(c fiber.Ctx) (uuid.UUID , error){
	userToken := jwtware.FromContext(c) 
	if userToken == nil {
		return uuid.Nil , errors.New("Cant find any usertoken.")
	}


	claims , claimsOk := userToken.Claims.(jwt.MapClaims)
	if !claimsOk {
		return  uuid.Nil , errors.New("Cant find claim in user token.")
	}

	uuidStr , uuidOk := claims["uuid"].(string)
	if !uuidOk {
		return uuid.Nil , errors.New("Cant convert uuid in token --> string")
	}

	parsedUUID, uuidErr := uuid.Parse(uuidStr)
	if uuidErr != nil {
		return  uuid.Nil , errors.New("Cant parse uuid string --> uuid.UUID")
	}

	return  parsedUUID , nil
}

func IsThisUserContainThisCane(c fiber.Ctx , cane_id string) (bool , error){
	userUUID , uuidErr := GetUUIDByContext(c)
	if uuidErr != nil {
		return  false , uuidErr
	}

	userCanes , canesErr := mygorm.GetCanesByUUID(userUUID)

	if canesErr != nil {
		return  false , canesErr
	}

	for _ , cane := range userCanes {
		if cane.CaneID == cane_id {
			return  true , nil
		}
	}

	return false, nil

}

func IsThisUserUUIDContainThisCane(uuid uuid.UUID , cane_id string) (bool , error){
	userCanes , canesErr := mygorm.GetCanesByUUID(uuid)

	if canesErr != nil {
		return  false , canesErr
	}

	for _ , cane := range userCanes {
		if cane.CaneID == cane_id {
			return  true , nil
		}
	}

	return false, nil

}


func AnyToUUID(input interface{}) (uuid.UUID, error) {
	if input == nil {
		return uuid.Nil, errors.New("input cannot be nil")
	}

	switch v := input.(type) {
	
	case uuid.UUID:
		return v, nil

	case string:
		return uuid.Parse(strings.TrimSpace(v))

	case []byte:
		return uuid.Parse(strings.TrimSpace(string(v)))

	case *uuid.UUID:
		if v == nil {
			return uuid.Nil, errors.New("input is a nil *uuid.UUID")
		}
		return *v, nil

	case *string:
		if v == nil {
			return uuid.Nil, errors.New("input is a nil *string")
		}
		return uuid.Parse(strings.TrimSpace(*v))

	default:
		return uuid.Nil, errors.New("unable to convert type" + v.(string) + "to uuid.UUID")
	}
}

