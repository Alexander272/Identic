package models

type SignIn struct {
	Username  string            `json:"username" binding:"required"`
	Password  string            `json:"password" binding:"required"`
	IPAddress string            `json:"-"`
	UserAgent string            `json:"-"`
	Metadata  *LoginMetadataDTO `json:"metadata,omitempty"`
}

type LoginEvent string

const (
	LoginEventSuccess                   LoginEvent = "login_success"
	LoginEventFailed                    LoginEvent = "login_failed"
	LoginEventSessionRefreshedAfterIdle LoginEvent = "session_refreshed_after_idle"
	LoginEventLogout                    LoginEvent = "logout"
)

type LoginMetadataDTO struct {
	Event          LoginEvent `json:"event,omitempty"`
	SessionID      string     `json:"sessionId,omitempty"`
	Geo            string     `json:"geo,omitempty"`
	City           string     `json:"city,omitempty"`
	Country        string     `json:"country,omitempty"`
	CountryCode    string     `json:"countryCode,omitempty"`
	Region         string     `json:"region,omitempty"`
	Device         string     `json:"device,omitempty"`
	OS             string     `json:"os,omitempty"`
	OSVersion      string     `json:"osVersion,omitempty"`
	Browser        string     `json:"browser,omitempty"`
	BrowserVersion string     `json:"browserVersion,omitempty"`
	IsMobile       bool       `json:"isMobile"`
	IsTablet       bool       `json:"isTablet"`
	IsDesktop      bool       `json:"isDesktop"`
	IsBot          bool       `json:"isBot"`
	Success        bool       `json:"success"`
	ErrorMessage   string     `json:"errorMessage,omitempty"`
}

type RefreshDTO struct {
	Token     string            `json:"token" binding:"required"`
	IPAddress string            `json:"-"`
	UserAgent string            `json:"-"`
	Metadata  *LoginMetadataDTO `json:"metadata,omitempty"`
}
