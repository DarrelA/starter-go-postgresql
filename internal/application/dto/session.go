package dto

// AuthSession contains the credentials issued to an authenticated browser.
type AuthSession struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}
