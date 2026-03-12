package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"shoego/config"
	"shoego/models"
)

func GetGoogleUser(code string) (*models.GoogleUser, error) {

	token, err := config.GoogleOAuthConfig.Exchange(
		context.Background(),
		code,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" +
			token.AccessToken,
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var user models.GoogleUser
	json.Unmarshal(body, &user)

	return &user, nil
}
