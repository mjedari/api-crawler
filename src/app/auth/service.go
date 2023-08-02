package auth

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mjedari/vgang-project/src/app/configs"
	"github.com/mjedari/vgang-project/src/domain/contracts"
	"io"
)

type AuthService struct {
	Client  *Client
	storage contracts.IStorage
	Token   string
}

func NewAuthService(storage contracts.IStorage) *AuthService {
	// todo: load from config
	client := NewClient("https://vgang.io/api/vgang-core/v1")

	return &AuthService{Client: client, storage: storage}
}

func (s *AuthService) setToken(token string) {
	s.Token = token
	s.Client.Token = token
}
func (s *AuthService) Login(ctx context.Context, request *LoginRequest) error {
	//check if it is in cache use it instead
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
	newPostRequest := PostRequest{
		Path:  "/auth/login/retailer/vgang",
		Body:  requestBody,
		Token: "",
	}
	res, err := s.Client.Post(context.Background(), newPostRequest)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)

	_ = json.Unmarshal(body, &authResponse)

	s.setToken(authResponse.GetAccessToken())

	err = s.storage.Store(ctx, "credential", string(body), 0)
	if err != nil {
		return err
	}

	return nil
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
