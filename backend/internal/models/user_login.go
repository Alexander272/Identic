package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserLogin struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	UserID         string          `json:"userId" db:"user_id"`
	LoginAt        time.Time       `json:"loginAt" db:"login_at"`
	IPAddress      *string         `json:"ipAddress" db:"ip_address"`
	UserAgent      *string         `json:"userAgent" db:"user_agent"`
	Metadata       json.RawMessage `json:"metadata" db:"metadata"`
	LastActivityAt time.Time       `json:"lastActivityAt" db:"last_activity_at"`
}

type UserLoginDTO struct {
	UserID    string          `json:"userId"`
	IPAddress *string         `json:"ipAddress"`
	UserAgent *string         `json:"userAgent"`
	Metadata  json.RawMessage `json:"metadata"`
}

type GetUserLoginsDTO struct {
	UserID    string     `json:"userId"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

type LoginMetadata struct {
	SessionID      string `json:"sessionId,omitempty"`
	Geo            string `json:"geo,omitempty"`
	City           string `json:"city,omitempty"`
	Country        string `json:"country,omitempty"`
	CountryCode    string `json:"countryCode,omitempty"`
	Region         string `json:"region,omitempty"`
	Device         string `json:"device,omitempty"`
	OS             string `json:"os,omitempty"`
	OSVersion      string `json:"osVersion,omitempty"`
	Browser        string `json:"browser,omitempty"`
	BrowserVersion string `json:"browserVersion,omitempty"`
	IsMobile       bool   `json:"isMobile"`
	IsTablet       bool   `json:"isTablet"`
	IsDesktop      bool   `json:"isDesktop"`
	IsBot          bool   `json:"isBot"`
	Success        bool   `json:"success"`
}
