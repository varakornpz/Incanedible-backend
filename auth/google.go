package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/varakornpz/mygorm"
	"github.com/varakornpz/models"
	"github.com/varakornpz/providers"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)



var googleOauthConfig *oauth2.Config

func InitGoogleAuthConf(){
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	type MyOAuthConfig struct {
		RedirectURL string
		ClientID	string
		ClientSecret	string
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  providers.AppConf.GGRedirectUrl,
		ClientID:     providers.AppConf.GGClientID,
		ClientSecret: providers.AppConf.GGClientSecret,                      
		Scopes:       []string{
			"https://www.googleapis.com/auth/userinfo.email", 
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleAuthCallBack(c fiber.Ctx) error{
	stateFromUrl := c.Query("state")
	stateFromCookie := c.Cookies("authstate")

	if stateFromCookie != stateFromUrl {
		return c.Status(fiber.ErrUnauthorized.Code).SendString("If you are not hacker(Please dont do anything to my shitty api) ,Please try again.")
	}

	code := c.Query("code") 
	if code == "" {
		return  c.Status(fiber.ErrUnauthorized.Code).SendString("This shit could not happend if you are not the hacker😭")
	}
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
	}
	res , err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer res.Body.Close()
	contents, err := io.ReadAll(res.Body)
	if err != nil {
			log.Error().Msg("Cant read user info from google api")
			return c.SendStatus(fiber.StatusInternalServerError)
	}
	var userInfo map[string]interface{}
	json.Unmarshal(contents, &userInfo)

	// email := userInfo["email"].(string)
	// name := userInfo["name"].(string)
	email := userInfo["email"].(string)
	name  := userInfo["name"].(string)

	var profileImage string
    if pic, ok := userInfo["picture"].(string); ok {
        profileImage = pic
    } else {
        profileImage = "https://your-domain.com/default-avatar.png" 
    }

	var currentUser models.User
	user , getUserErr := mygorm.GetUserByEmail(email)
	if getUserErr != nil {
		currentUser.Email = email
		currentUser.Name  = name
		currentUser.ProfilePic = profileImage
		mygorm.PutNewUser(&currentUser)
	}else{
		currentUser = user
	}

	jwtToken , tokennErr := getJWT(currentUser.UUID , providers.AppConf.JWTSecret)
	if tokennErr != nil {
		c.Redirect().Status(fiber.StatusFound).To(providers.AppConf.FrontendErrPage + "/?error=token_create_fail" )
	}


	log.Debug().Msg(currentUser.Email)
	c.Cookie(&fiber.Cookie{
			Name:     	"access_token",
			Value:		jwtToken,
			//Expired next 24hr
			Expires:  	time.Now().Add(24 * time.Hour),
			//Pretend frontend to edit cookie
			HTTPOnly: 	true,
			//***In prod change to true***
			Secure:  	 true,
			SameSite: 	"Lax",
			Path:     	"/",
			Domain: providers.AppConf.COOKIEDomain,
		})
	return c.Redirect().Status(fiber.StatusTemporaryRedirect).To(providers.AppConf.GGAfterSigninRedirect)
}

func GoogleAuthSignin(c fiber.Ctx) error{
	// if c.Cookies("access_token") != ""{
	// 	claims ,err := utils.Ver
	// }
	oAuthState := genCSRFToken(c)

	url := googleOauthConfig.AuthCodeURL(oAuthState)

	return  c.Redirect().Status(fiber.StatusTemporaryRedirect).To(url)
}