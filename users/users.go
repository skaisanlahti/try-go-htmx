package users

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/logging"
	"github.com/skaisanlahti/try-go-htmx/users/handlers"
	"github.com/skaisanlahti/try-go-htmx/users/passwords"
	"github.com/skaisanlahti/try-go-htmx/users/repositories"
	"github.com/skaisanlahti/try-go-htmx/users/sessions"
	"github.com/skaisanlahti/try-go-htmx/users/templates"
)

func UseUserRoutes(router *http.ServeMux, database *sql.DB, store *sessions.Store) {
	htmlTemplates := templates.ParseTemplates()
	userRepository := repositories.NewUserRepository(database)
	passwordEncoder := passwords.NewArgon2Encoder(passwords.DefaultArgon2idOptions)
	// passwordEncoder := passwords.NewBcryptEncoder(12)

	getRegisterPage := handlers.NewGetRegisterPageHandler(handlers.NewHtmxGetRegisterPageRenderer(htmlTemplates.RegisterPage))
	getLoginPage := handlers.NewGetLoginPageHandler(handlers.NewHtmxGetLoginPageRenderer(htmlTemplates.LoginPage))
	getLogoutPage := handlers.NewGetLogoutPageHandler(handlers.NewHtmxGetLogoutPageRenderer(htmlTemplates.LogoutPage))
	addUser := handlers.NewAddUserHandler(passwordEncoder, userRepository, store, handlers.NewHtmxAddUserRenderer(htmlTemplates.RegisterPage))
	logoutUser := handlers.NewLogoutUserHandler(store)
	loginUser := handlers.NewLoginUserHandler(passwordEncoder, userRepository, store, handlers.NewHtmxLoginUserRender(htmlTemplates.LoginPage),
		handlers.LoginUserOptions{RecalculateOutdatedKeys: true},
	)

	router.Handle("/users/logout", logging.LogRequest(sessions.RequireSession(logoutUser, store)))
	router.Handle("/users/login", logging.LogRequest(loginUser))
	router.Handle("/users/add", logging.LogRequest(addUser))
	router.Handle("/logout", logging.LogRequest(getLogoutPage))
	router.Handle("/login", logging.LogRequest(getLoginPage))
	router.Handle("/register", logging.LogRequest(getRegisterPage))
}
