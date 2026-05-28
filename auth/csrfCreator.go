package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v3"
	//"github.com/varakornpz/providers"
)

func genCSRFToken(c fiber.Ctx) string {
	byteFill := make([]byte , 16)
	rand.Read(byteFill)
	state := base64.URLEncoding.EncodeToString(byteFill)

	c.Cookie(&fiber.Cookie{
			Name:     "authstate",
			Value:    state,
			Expires:  time.Now().Add(10 * time.Minute),
			Path:     "/", 
        	Domain:   ".varakorn.net",
			HTTPOnly: true,
			SameSite: "Lax",
			Secure:   false,
		})

		return state
}