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

type AddUserSessionManager interface {
	Add(userId int) (*http.Cookie, error)
}

type AddUserHandler struct {
	repository AddUserRepository
	sessions   AddUserSessionManager
	renderer   AddUserRenderer
}

func NewAddUserHandler(
	repository AddUserRepository,
	sessions AddUserSessionManager,
	renderer AddUserRenderer,
) *AddUserHandler {
	return &AddUserHandler{repository, sessions, renderer}
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
		renderError("User name is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
	}

	user, err := handler.repository.GetUserByName(name)
	if err == nil {
		renderError("User name already exists.")
		return
	}

	user, err = domain.NewUser(name, password)
	if err != nil {
		renderError(err.Error())
		return
	}

	err = handler.repository.AddUser(user)
	if err != nil {
		renderError(err.Error())
		return
	}

	addedUser, _ := handler.repository.GetUserByName(user.Name)
	cookie, err := handler.sessions.Add(addedUser.Id)
	if err != nil {
		renderError(err.Error())
		return
	}

	http.SetCookie(response, cookie)
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
