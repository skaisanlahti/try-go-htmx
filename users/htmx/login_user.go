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
}

func NewLoginUserHandler(
	repository LoginUserRepository,
	sessions LoginUserSessionStore,
	renderer LoginUserRenderer,
) *LoginUserHandler {
	return &LoginUserHandler{repository, sessions, renderer}
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

	user, err := handler.repository.GetUserByName(name)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	isPasswordValid := domain.IsPasswordValid(user.Password, []byte(password))
	if err != nil || !isPasswordValid {
		renderError("Invalid credentials.")
		return
	}

	cookie, err := handler.sessions.Add(user.Id)
	if err != nil {
		renderError(err.Error())
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
