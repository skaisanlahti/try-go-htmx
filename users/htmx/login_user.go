package htmx

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type LoginUserRepository interface {
	GetUserByName(name string) (domain.User, error)
}

type LoginUserSessionStore interface {
	Add(userId int) (*http.Cookie, error)
}

type LoginUserRenderer interface {
	RenderLoginForm(name string, password string, errorMessage string) []byte
}

type LoginUserHandler struct {
	repository LoginUserRepository
	sessions   LoginUserSessionStore
	renderer   LoginUserRenderer
	fakeUser   domain.User
}

func NewLoginUserHandler(
	repository LoginUserRepository,
	sessions LoginUserSessionStore,
	renderer LoginUserRenderer,
) *LoginUserHandler {
	fakeUser, err := domain.NewUser("FakeName", "Fake Passphrase to compare")
	if err != nil {
		log.Panicln("Failed to create fake user for login.")
	}

	return &LoginUserHandler{repository, sessions, renderer, fakeUser}
}

func (handler *LoginUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := handler.renderer.RenderLoginForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	}

	isPasswordValid := false
	user, err := handler.repository.GetUserByName(name)
	if err != nil {
		domain.IsPasswordValid(handler.fakeUser.Password, []byte(password))
	} else {
		isPasswordValid = domain.IsPasswordValid(user.Password, []byte(password))
	}

	if !isPasswordValid {
		renderError("Invalid credentials.")
		return
	}

	cookie, err := handler.sessions.Add(user.Id)
	if err != nil {
		renderError("Internal session error.")
		return
	}

	http.SetCookie(response, cookie)
	response.Header().Add("HX-Redirect", "/todos")
	response.WriteHeader(http.StatusOK)
}

type HtmxLoginUserRenderer struct {
	loginPage *template.Template
}

func NewHtmxLoginUserRender(loginPage *template.Template) *HtmxLoginUserRenderer {
	return &HtmxLoginUserRenderer{loginPage}
}

func (renderer *HtmxLoginUserRenderer) RenderLoginForm(name string, password string, errorMessage string) []byte {
	templateData := LoginPage{time.Now().UnixMilli(), name, password, errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.ExecuteTemplate(buffer, "form", templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
