package htmx

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type AddUserRepository interface {
	GetUserByName(name string) (domain.User, error)
	AddUser(user domain.User) error
}

type AddUserRenderer interface {
	RenderAddUserForm(name string, password string, errorMessage string) []byte
}

type AddUserSessionStore interface {
	Add(userId int) *domain.Session
}

type AddUserHandler struct {
	repository AddUserRepository
	sessions   AddUserSessionStore
	renderer   AddUserRenderer
	mode       string
}

func NewAddUserHandler(
	repository AddUserRepository,
	sessions AddUserSessionStore,
	renderer AddUserRenderer,
	mode string,
) *AddUserHandler {
	return &AddUserHandler{repository, sessions, renderer, mode}
}

func (handler *AddUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	if name == "" {
		html := handler.renderer.RenderAddUserForm(name, password, "User name is required.")
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	if password == "" {
		html := handler.renderer.RenderAddUserForm(name, password, "Password is required.")
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	user, err := handler.repository.GetUserByName(name)
	if err == nil {
		html := handler.renderer.RenderAddUserForm(name, password, "User name already exists.")
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	user, err = domain.NewUser(name, password)
	if err != nil {
		html := handler.renderer.RenderAddUserForm(name, password, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	err = handler.repository.AddUser(user)
	if err != nil {
		html := handler.renderer.RenderAddUserForm(name, password, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	addedUser, err := handler.repository.GetUserByName(user.Name)
	session := handler.sessions.Add(addedUser.Id)
	http.SetCookie(response, domain.NewSessionCookie(session, handler.mode))
	response.Header().Add("HX-Redirect", "/todos")
	response.WriteHeader(http.StatusOK)
}

type HtmxAddUserRenderer struct {
	addUserPage *template.Template
}

func NewHtmxAddUserRenderer(addUserPage *template.Template) *HtmxAddUserRenderer {
	return &HtmxAddUserRenderer{addUserPage}
}

func (renderer *HtmxAddUserRenderer) RenderAddUserForm(name string, password string, errorMessage string) []byte {
	templateData := RegisterPage{time.Now().UnixMilli(), name, password, errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.addUserPage.ExecuteTemplate(buffer, "form", templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
