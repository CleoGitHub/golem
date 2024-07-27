package httpclient

type SecurityClient interface {
	Authenticate(email string, password string) (token string, refreshToken string, err error)
	RefreshToken(refresh string) (token string, refreshToken string, err error)
}
