package users

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/middleware"
	"github.com/skaisanlahti/try-go-htmx/sessions"
	"github.com/skaisanlahti/try-go-htmx/users/htmx"
	"github.com/skaisanlahti/try-go-htmx/users/psql"
	"github.com/skaisanlahti/try-go-htmx/users/templates"
)

type SessionStore interface {
	Add(userId int) (*http.Cookie, error)
	Remove(*http.Request) (*http.Cookie, error)
	Validate(*http.Request) (*sessions.Session, error)
	Extend(*sessions.Session) (*http.Cookie, error)
}

func MapHtmxHandlers(router *http.ServeMux, database *sql.DB, sessions SessionStore) {
	htmlTemplates := templates.ParseTemplates()
	getLoginPage := htmx.NewGetLoginPageHandler(htmx.NewHtmxGetLoginPageRenderer(htmlTemplates.LoginPage))
	getLogoutPage := htmx.NewGetLogoutPageHandler(htmx.NewHtmxGetLogoutPageRenderer(htmlTemplates.LogoutPage))
	getRegisterPage := htmx.NewGetRegisterPageHandler(htmx.NewHtmxGetRegisterPageRenderer(htmlTemplates.RegisterPage))

	userRepository := psql.NewUserRepository(database)
	loginUser := htmx.NewLoginUserHandler(userRepository, sessions, htmx.NewHtmxLoginUserRender(htmlTemplates.LoginPage))
	logoutUser := htmx.NewLogoutUserHandler(sessions)
	addUser := htmx.NewAddUserHandler(userRepository, sessions, htmx.NewHtmxAddUserRenderer(htmlTemplates.RegisterPage))

	router.Handle("/users/logout", middleware.LogRequest(middleware.RequireSession(logoutUser, sessions)))
	router.Handle("/users/login", middleware.LogRequest(loginUser))
	router.Handle("/users/add", middleware.LogRequest(addUser))
	router.Handle("/logout", middleware.LogRequest(getLogoutPage))
	router.Handle("/login", middleware.LogRequest(getLoginPage))
	router.Handle("/register", middleware.LogRequest(getRegisterPage))
}
