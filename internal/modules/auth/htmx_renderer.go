package auth

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"time"
)

type htmxRenderer struct {
	loginPage    *template.Template
	logoutPage   *template.Template
	registerPage *template.Template
}

func newHtmxRenderer(files embed.FS) *htmxRenderer {
	loginPage := template.Must(template.ParseFS(files, "web/html/login_page.html"))
	logoutPage := template.Must(template.ParseFS(files, "web/html/logout_page.html"))
	registerPage := template.Must(template.ParseFS(files, "web/html/register_page.html"))
	return &htmxRenderer{loginPage, logoutPage, registerPage}
}

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (renderer *htmxRenderer) renderRegisterPage() []byte {
	data := registerPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.registerPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmxRenderer) renderRegisterForm(name string, password string, errorMessage string) []byte {
	data := registerPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.registerPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type loginPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (renderer *htmxRenderer) renderLoginPage() []byte {
	data := loginPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmxRenderer) renderLoginForm(name string, password string, errorMessage string) []byte {
	data := loginPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type logoutPageData struct {
	LoggedIn bool
}

func (renderer *htmxRenderer) renderLogoutPage(loggedIn bool) []byte {
	data := logoutPageData{loggedIn}
	buffer := &bytes.Buffer{}
	err := renderer.logoutPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}