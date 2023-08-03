package auth

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mjedari/vgang-project/app/configs"
	"github.com/mjedari/vgang-project/domain/contracts"
	"github.com/mjedari/vgang-project/infra/client"
	"io"
	"time"
)

type AuthService struct {
	Client  contracts.IHTTPClient
	storage contracts.IStorage
	Token   string
	config  configs.OriginRemote
}

func NewAuthService(storage contracts.IStorage, config configs.OriginRemote) *AuthService {
	client := client.NewClient(config.BaseURL)
	return &AuthService{storage: storage, Client: client, config: config}
}

func (s *AuthService) Login(ctx context.Context, request *LoginRequest) error {
	// todo: ctx
	var authResponse AuthResponse

	data := s.storage.Fetch(ctx, "credential")
	err := json.Unmarshal(data, &authResponse)
	// todo: handle err
	// todo: refresh token

	token := authResponse.GetAccessToken()

	if token != "" {
		s.setToken(token)
		return nil
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	newPostRequest := client.PostRequest{
		Path:  s.config.Login,
		Body:  requestBody,
		Token: "",
	}
	res, err := s.Client.Post(ctx, newPostRequest)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)

	_ = json.Unmarshal(body, &authResponse)

	s.setToken(authResponse.GetAccessToken())

	err = s.storage.Store(ctx, "credential", string(body), time.Hour)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) setToken(token string) {
	s.Token = token
	s.Client.SetToken(token)
}

type AccessTokens struct {
	Token   string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

type AuthResponse struct {
	Data AccessTokens `json:"data"`
}

func (r *AuthResponse) GetAccessToken() string {
	return r.Data.Token
}

type LoginRequest struct {
	DeviceId string
	Email    string
	Password string
}

func NewLoginRequest() *LoginRequest {
	// todo: make it done in wiring section
	deviceId, _ := uuid.NewUUID()
	userName := configs.Config.Credentials.UserName
	password := configs.Config.Credentials.Password
	return &LoginRequest{DeviceId: deviceId.String(), Email: userName, Password: password}
}

type LoginResponse struct {
	Token string
}
