package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/google/uuid"
)

type SessionService struct {
	keycloak *auth.KeycloakClient
	user     Users
	policies AccessPolices
}

func NewSessionService(keycloak *auth.KeycloakClient, user Users, policies AccessPolices) *SessionService {
	return &SessionService{
		keycloak: keycloak,
		user:     user,
		policies: policies,
	}
}

type Session interface {
	SignIn(ctx context.Context, u *models.SignIn) (*models.User, error)
	SignOut(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, req *models.RefreshDTO) (*models.User, error)
	DecodeAccessToken(ctx context.Context, token string) (*models.User, error)
}

func (s *SessionService) SignIn(ctx context.Context, u *models.SignIn) (*models.User, error) {
	res, err := s.keycloak.Client.Login(ctx, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm, u.Username, u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login to keycloak. error: %w", err)
	}

	decodedUser, err := s.DecodeAccessToken(ctx, res.AccessToken)
	if err != nil {
		return nil, err
	}
	access, err := s.policies.GetPolicies(decodedUser.ID.String())
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           decodedUser.ID,
		Name:         decodedUser.Name,
		Role:         access.Role,
		Permissions:  access.Perms,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	return user, nil
}

func (s *SessionService) SignOut(ctx context.Context, refreshToken string) error {
	err := s.keycloak.Client.Logout(ctx, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to logout to keycloak. error: %w", err)
	}
	return nil
}

func (s *SessionService) Refresh(ctx context.Context, req *models.RefreshDTO) (*models.User, error) {
	res, err := s.keycloak.Client.RefreshToken(ctx, req.Token, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token in keycloak. error: %w", err)
	}

	decodedUser, err := s.DecodeAccessToken(ctx, res.AccessToken)
	if err != nil {
		return nil, err
	}
	access, err := s.policies.GetPolicies(decodedUser.ID.String())
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           decodedUser.ID,
		Name:         decodedUser.Name,
		Role:         access.Role,
		Permissions:  access.Perms,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	return user, nil
}

func (s *SessionService) DecodeAccessToken(ctx context.Context, token string) (*models.User, error) {
	//TODO расшифровку токена тоже лучше делать здесь, а в keycloak
	_, claims, err := s.keycloak.Client.DecodeAccessToken(ctx, token, s.keycloak.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode access token. error: %w", err)
	}

	serviceName := os.Getenv("SERVICE_ID")

	user := &models.User{}
	var role, username string
	var userId uuid.UUID
	c := *claims
	access, ok := c["realm_access"]
	if ok {
		a := access.(map[string]interface{})["roles"]
		roles := a.([]interface{})
		for _, r := range roles {
			//TODO может получать прификс из конфига
			if strings.Contains(r.(string), serviceName) {
				role = strings.Replace(r.(string), serviceName+"_", "", 1)
				break
			}
		}
	}

	u, ok := c["preferred_username"]
	if ok {
		username = u.(string)
	}

	uId, ok := c["sub"]
	if ok {
		strId := uId.(string)
		userId, err = uuid.Parse(strId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user id. error: %w", err)
		}
	}

	user.ID = userId
	user.Role = role
	user.Name = username

	return user, nil
}
