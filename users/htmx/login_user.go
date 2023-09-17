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
	Add(userId int) *domain.Session
}

type LoginUserRenderer interface {
	RenderLoginForm(name string, password string, errorMessage string) []byte
}

type LoginUserHandler struct {
	repository LoginUserRepository
	sessions   LoginUserSessionStore
	renderer   LoginUserRenderer
	mode       string
}

func NewLoginUserHandler(
	repository LoginUserRepository,
	sessions LoginUserSessionStore,
	renderer LoginUserRenderer,
	mode string,
) *LoginUserHandler {
	return &LoginUserHandler{repository, sessions, renderer, mode}
}

func (handler *LoginUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	log.Printf("Name: %s Password: %s", name, password)

	user, err := handler.repository.GetUserByName(name)
	if err != nil {
		html := handler.renderer.RenderLoginForm(name, password, "Invalid credentials.")
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	isPasswordValid := domain.IsPasswordValid(user.Password, []byte(password))
	if err != nil || !isPasswordValid {
		html := handler.renderer.RenderLoginForm(name, password, "Invalid credentials.")
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	session := handler.sessions.Add(user.Id)
	http.SetCookie(response, domain.NewSessionCookie(session, handler.mode))
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
