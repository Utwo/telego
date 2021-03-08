package middlewares

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"strings"
	"telego/app/models"
	"telego/app/services"
)

var Auth = func(isOptional bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" && isOptional {
				return next(c)
			}
			if authHeader == "" && !isOptional {
				return echo.ErrForbidden
			}
			jwtToken := strings.Fields(authHeader)[1]
			if jwtToken == "" {
				return echo.ErrForbidden
			}

			firebase := c.Get("firebase").(*auth.Client)
			ctx := c.Get("firebaseCtx").(context.Context)
			token, err := firebase.VerifyIDToken(ctx, jwtToken)
			if err != nil {
				return echo.ErrForbidden
			}

			log.Printf("Verified ID token: %v\n", token)
			db := c.Get("db").(*gorm.DB)
			account := models.Account{
				AuthId: token.UID,
				//Name:                  token.Firebase.Identities,
				//Email:                 "",
				//Picture:               "",
			}
			tx := db.Where(models.Account{AuthId: account.AuthId}).FirstOrCreate(&account)
			if tx.RowsAffected > 0 && !account.IsAnonymous {
				//new user
				// TODO: handleActiveInvitations
				services.SendWelcomeMail(account.Email, account.Name)
			}
			c.Set("account", account)
			return next(c)
		}
	}
}
