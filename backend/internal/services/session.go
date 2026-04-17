package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/google/uuid"
)

type SessionService struct {
	keycloak   *auth.KeycloakClient
	user       Users
	policies   AccessPolices
	userLogins UserLogins
}

func NewSessionService(keycloak *auth.KeycloakClient, user Users, policies AccessPolices, userLogins UserLogins) *SessionService {
	return &SessionService{
		keycloak:   keycloak,
		user:       user,
		policies:   policies,
		userLogins: userLogins,
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
		s.recordFailedLogin(context.Background(), u, err.Error())
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

	go func() {
		s.recordSuccessfulLogin(context.Background(), decodedUser.ID.String(), u)
	}()

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

	go func() {
		s.checkAndRecordIdleSession(context.Background(), decodedUser.ID.String(), req)
	}()

	return user, nil
}

func (s *SessionService) recordSuccessfulLogin(ctx context.Context, userID string, u *models.SignIn) {
	metadata := ParseUserAgent(u.UserAgent)
	metadata.Event = models.LoginEventSuccess
	metadata.Success = true

	if u.Metadata != nil {
		metadata.SessionID = u.Metadata.SessionID
		metadata.Geo = u.Metadata.Geo
		metadata.City = u.Metadata.City
		metadata.Country = u.Metadata.Country
		metadata.CountryCode = u.Metadata.CountryCode
		metadata.Region = u.Metadata.Region
	}

	data, _ := json.Marshal(metadata)

	ip := u.IPAddress
	ua := u.UserAgent

	dto := &models.UserLoginDTO{
		UserID:    userID,
		IPAddress: &ip,
		UserAgent: &ua,
		Metadata:  data,
	}
	if err := s.userLogins.RecordLogin(ctx, dto); err != nil {
		logger.Error("failed to record successful login", logger.ErrAttr(err))
		error_bot.Send(nil, fmt.Sprintf("failed to record successful login. error: %v", err), dto)
	}
}

func (s *SessionService) recordFailedLogin(ctx context.Context, u *models.SignIn, errMsg string) {
	metadata := ParseUserAgent(u.UserAgent)
	metadata.Event = models.LoginEventFailed
	metadata.Success = false
	metadata.ErrorMessage = errMsg

	if u.Metadata != nil {
		metadata.SessionID = u.Metadata.SessionID
		metadata.Geo = u.Metadata.Geo
		metadata.City = u.Metadata.City
		metadata.Country = u.Metadata.Country
		metadata.CountryCode = u.Metadata.CountryCode
		metadata.Region = u.Metadata.Region
	}

	data, _ := json.Marshal(metadata)

	ip := u.IPAddress
	ua := u.UserAgent

	dto := &models.UserLoginDTO{
		UserID:    "",
		IPAddress: &ip,
		UserAgent: &ua,
		Metadata:  data,
	}
	if err := s.userLogins.RecordLogin(ctx, dto); err != nil {
		logger.Error("failed to record failed login", logger.ErrAttr(err))
		error_bot.Send(nil, fmt.Sprintf("failed to record failed login. error: %v", err), dto)
	}
}

func (s *SessionService) checkAndRecordIdleSession(ctx context.Context, userID string, req *models.RefreshDTO) {
	wasIdle, err := s.userLogins.UpdateLastActivity(ctx, nil, userID)
	if err != nil {
		logger.Error("failed to update last activity", logger.ErrAttr(err))
		error_bot.Send(nil, fmt.Sprintf("failed to update last activity. error: %v", err), userID)
		return
	}
	if !wasIdle {
		return
	}

	metadata := ParseUserAgent(req.UserAgent)
	metadata.Event = models.LoginEventSessionRefreshedAfterIdle
	metadata.Success = true

	if req.Metadata != nil {
		metadata.Geo = req.Metadata.Geo
		metadata.City = req.Metadata.City
		metadata.Country = req.Metadata.Country
		metadata.CountryCode = req.Metadata.CountryCode
		metadata.Region = req.Metadata.Region
	}

	data, _ := json.Marshal(metadata)

	ip := req.IPAddress
	ua := req.UserAgent

	if err := s.userLogins.RecordLogin(ctx, &models.UserLoginDTO{
		UserID:    userID,
		IPAddress: &ip,
		UserAgent: &ua,
		Metadata:  data,
	}); err != nil {
		logger.Error("failed to record idle session", logger.ErrAttr(err))
		error_bot.Send(nil, "failed to record idle session", err)
	}
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
