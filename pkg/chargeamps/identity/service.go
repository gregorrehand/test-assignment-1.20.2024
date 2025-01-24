package identity

import (
	"context"
	"encoding/json"
	"gitlab.com/gridio/test-assignment/pkg/chargeamps/utils"

	"github.com/sirupsen/logrus"
	"gitlab.com/gridio/test-assignment/internal"
)

type TokenSource struct {
	token internalToken
}

type internalToken struct {
	// TODO: Replace with proper fields that chargeamps is responding with
	Message      string `json:"message"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
		RfidTags  *[]struct {
			Active         bool   `json:"active"`
			Rfid           string `json:"rfid"`
			RfidDec        string `json:"rfidDec"`
			RfidDecReverse string `json:"rfidDecReverse"`
		} `json:"rfidTags"`
		UserStatus string `json:"userStatus"`
	} `json:"user"`
}

type refreshTokenPayload struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateFromSecretAgent(logger logrus.FieldLogger, sa internal.SecretAgent) *TokenSource {
	t := TokenSource{}

	var unmarshalled internalToken

	err := json.Unmarshal([]byte(sa.ProvideSecret()), &unmarshalled)
	// TODO: Check error here
	if err != nil {
		logger.Error("Failed to unmarshal secret agent data: ", err)
		return nil
	}

	t.token = unmarshalled

	return &t
}

func Login(ctx context.Context, logger logrus.FieldLogger, apiClient *utils.APIClient, username, password string) (*TokenSource, error) {
	t := TokenSource{}

	payload := loginPayload{
		Email:    username,
		Password: password,
	}

	err := apiClient.PostWithoutToken(ctx, "auth/login", payload, &t.token)
	if err != nil {
		logger.Error("Failed to log in: ", err)
		return nil, err
	}

	return &t, nil
}

func (t *TokenSource) AccessToken() string {
	// TODO implement me
	return t.token.Token
}

func (t *TokenSource) IsUnauthorized() bool {
	// TODO implement me
	return t.token.Token == ""
}

func (t *TokenSource) String() string {
	b, _ := json.Marshal(t.token)
	return string(b)
}

// TODO: Write a function that retrieves access and refresh tokens from chargeamps and stores them in internalToken
// 	struct

func getRefreshToken(ctx context.Context, logger logrus.FieldLogger, apiClient *utils.APIClient, existingToken, refreshToken string) (*internalToken, error) {
	payload := refreshTokenPayload{
		Token:        existingToken,
		RefreshToken: refreshToken,
	}

	var token internalToken

	err := apiClient.PostWithoutToken(ctx, "auth/refreshToken", payload, &token)
	if err != nil {
		logger.Error("Failed to fetch refresh token: ", err)
		return nil, err
	}

	return &token, nil
}
