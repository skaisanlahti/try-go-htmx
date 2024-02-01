package security

import "net/http"

type CookieFactory struct {
	options SessionOptions
}

func NewCookieFactory(options SessionOptions) *CookieFactory {
	return &CookieFactory{options}
}

func (this *CookieFactory) NewSessionCookie(signedSession string) *http.Cookie {
	return &http.Cookie{
		Name:     this.options.CookieName,
		Path:     "/",
		Value:    signedSession,
		MaxAge:   int(this.options.Duration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   this.options.Secure,
	}
}

func (this *CookieFactory) ClearSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     this.options.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   this.options.Secure,
	}
}
