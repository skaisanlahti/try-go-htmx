package users

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/middleware"
	"github.com/skaisanlahti/try-go-htmx/users/domain"
	"github.com/skaisanlahti/try-go-htmx/users/htmx"
	"github.com/skaisanlahti/try-go-htmx/users/psql"
	"github.com/skaisanlahti/try-go-htmx/users/templates"
)

type SessionStore interface {
	Add(userId int) *domain.Session
	Remove(sessionId string)
	Validate(sessionId string) (*domain.Session, bool)
	Extend(*domain.Session) *domain.Session
}

func MapHtmxHandlers(router *http.ServeMux, database *sql.DB, sessions SessionStore, mode string) {
	htmlTemplates := templates.ParseTemplates()
	getLoginPage := htmx.NewGetLoginPageHandler(htmx.NewHtmxGetLoginPageRenderer(htmlTemplates.LoginPage))
	getLogoutPage := htmx.NewGetLogoutPageHandler(htmx.NewHtmxGetLogoutPageRenderer(htmlTemplates.LogoutPage))
	getRegisterPage := htmx.NewGetRegisterPageHandler(htmx.NewHtmxGetRegisterPageRenderer(htmlTemplates.RegisterPage))

	userRepository := psql.NewUserRepository(database)
	loginUser := htmx.NewLoginUserHandler(userRepository, sessions, htmx.NewHtmxLoginUserRender(htmlTemplates.LoginPage), mode)
	logoutUser := htmx.NewLogoutUserHandler(sessions)
	addUser := htmx.NewAddUserHandler(userRepository, sessions, htmx.NewHtmxAddUserRenderer(htmlTemplates.RegisterPage), mode)

	router.Handle("/users/logout", middleware.LogRequest(middleware.RequireSession(logoutUser, sessions, mode)))
	router.Handle("/users/login", middleware.LogRequest(loginUser))
	router.Handle("/users/add", middleware.LogRequest(addUser))
	router.Handle("/logout", middleware.LogRequest(getLogoutPage))
	router.Handle("/login", middleware.LogRequest(getLoginPage))
	router.Handle("/register", middleware.LogRequest(getRegisterPage))
}
