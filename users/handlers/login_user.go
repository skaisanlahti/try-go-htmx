package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type LoginUserPasswordEncoder interface {
	NewKey(password string) ([]byte, error)
	VerifyKey(encodedKey []byte, candidatePassword string, recalculateOutdatedKeys bool) (bool, chan []byte)
}

type LoginUserRepository interface {
	GetUserByName(name string) (domain.User, error)
	UpdateUserPassword(user domain.User) error
}

type LoginUserSession interface {
	StartSession(response http.ResponseWriter, userId int) error
}

type LoginUserRenderer interface {
	RenderLoginForm(name string, password string, errorMessage string) []byte
}

type LoginUserOptions struct {
	RecalculateOutdatedKeys bool
}

type LoginUserHandler struct {
	encoder    LoginUserPasswordEncoder
	repository LoginUserRepository
	session    LoginUserSession
	renderer   LoginUserRenderer
	options    LoginUserOptions
	fakeUser   domain.User
}

func NewLoginUserHandler(
	encoder LoginUserPasswordEncoder,
	repository LoginUserRepository,
	session LoginUserSession,
	renderer LoginUserRenderer,
	options LoginUserOptions,
) *LoginUserHandler {
	fakeKey, err := encoder.NewKey("Fake password to compare")
	if err != nil {
		log.Panicln("Failed to create fake key for login.")
	}

	fakeUser, err := domain.NewUser("FakeName", fakeKey)
	if err != nil {
		log.Panicln("Failed to create fake user for login.")
	}

	return &LoginUserHandler{encoder, repository, session, renderer, options, fakeUser}
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
		handler.encoder.VerifyKey(handler.fakeUser.Password, password, handler.options.RecalculateOutdatedKeys)
		renderError("Invalid credentials.")
		return
	}

	isPasswordCorrect, newKeyChannel := handler.encoder.VerifyKey(user.Password, password, handler.options.RecalculateOutdatedKeys)
	if !isPasswordCorrect {
		renderError("Invalid credentials.")
		return
	}

	if newKeyChannel != nil {
		go handler.updatePassword(user, newKeyChannel)
	}

	err = handler.session.StartSession(response, user.Id)
	if err != nil {
		renderError("Internal session error.")
		return
	}

	response.Header().Add("HX-Location", "/todos")
	response.WriteHeader(http.StatusOK)
}

func (handler *LoginUserHandler) updatePassword(user domain.User, newKeyResult <-chan []byte) {
	newKey, ok := <-newKeyResult
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Password = newKey
	err := handler.repository.UpdateUserPassword(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}

type HtmxLoginUserView struct {
	loginPage *template.Template
}

func NewHtmxLoginUserView(loginPage *template.Template) *HtmxLoginUserView {
	return &HtmxLoginUserView{loginPage}
}

func (renderer *HtmxLoginUserView) RenderLoginForm(name string, password string, errorMessage string) []byte {
	templateData := LoginPage{time.Now().UnixMilli(), name, password, errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.ExecuteTemplate(buffer, "form", templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
