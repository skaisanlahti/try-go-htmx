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

	logoutUser := handlers.NewLogoutUserHandler(store)
	router.Handle("/users/logout", logging.LogRequest(sessions.RequireSession(logoutUser, store)))

	loginUserView := handlers.NewHtmxLoginUserView(htmlTemplates.LoginPage)
	loginOptions := handlers.LoginUserOptions{RecalculateOutdatedKeys: true}
	loginUser := handlers.NewLoginUserHandler(passwordEncoder, userRepository, store, loginUserView, loginOptions)
	router.Handle("/users/login", logging.LogRequest(loginUser))

	addUserView := handlers.NewHtmxAddUserView(htmlTemplates.RegisterPage)
	addUser := handlers.NewAddUserHandler(passwordEncoder, userRepository, store, addUserView)
	router.Handle("/users/add", logging.LogRequest(addUser))

	getLogoutPageView := handlers.NewHtmxGetLogoutPageView(htmlTemplates.LogoutPage)
	getLogoutPage := handlers.NewGetLogoutPageHandler(getLogoutPageView)
	router.Handle("/logout", logging.LogRequest(getLogoutPage))

	getLoginPageView := handlers.NewHtmxGetLoginPageView(htmlTemplates.LoginPage)
	getLoginPage := handlers.NewGetLoginPageHandler(getLoginPageView)
	router.Handle("/login", logging.LogRequest(getLoginPage))

	getRegisterPageView := handlers.NewHtmxGetRegisterPageView(htmlTemplates.RegisterPage)
	getRegisterPage := handlers.NewGetRegisterPageHandler(getRegisterPageView)
	router.Handle("/register", logging.LogRequest(getRegisterPage))
}
