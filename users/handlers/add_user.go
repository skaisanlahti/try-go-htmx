package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type AddUserPasswordEncoder interface {
	NewKey(password string) ([]byte, error)
}

type AddUserRepository interface {
	GetUserByName(name string) (domain.User, error)
	AddUser(user domain.User) error
}

type AddUserRenderer interface {
	RenderAddUserForm(name string, password string, errorMessage string) []byte
}

type AddUserSessionStore interface {
	Add(userId int) (*http.Cookie, error)
}

type AddUserHandler struct {
	encoder    AddUserPasswordEncoder
	repository AddUserRepository
	sessions   AddUserSessionStore
	renderer   AddUserRenderer
}

func NewAddUserHandler(
	encoder AddUserPasswordEncoder,
	repository AddUserRepository,
	sessions AddUserSessionStore,
	renderer AddUserRenderer,
) *AddUserHandler {
	return &AddUserHandler{encoder, repository, sessions, renderer}
}

func (handler *AddUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := handler.renderer.RenderAddUserForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
	}

	newUser, err := handler.repository.GetUserByName(name)
	if err == nil {
		renderError("Username already exists.")
		return
	}

	key, err := handler.encoder.NewKey(password)
	if err != nil {
		renderError(err.Error())
		return
	}

	newUser, err = domain.NewUser(name, key)
	if err != nil {
		renderError(err.Error())
		return
	}

	err = handler.repository.AddUser(newUser)
	if err != nil {
		renderError(err.Error())
		return
	}

	newUser, _ = handler.repository.GetUserByName(newUser.Name)
	cookie, err := handler.sessions.Add(newUser.Id)
	if err != nil {
		renderError(err.Error())
		return
	}

	http.SetCookie(response, cookie)
	response.Header().Add("HX-Location", "/todos")
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
