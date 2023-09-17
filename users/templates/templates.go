package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var templateFiles embed.FS

type HtmlTemplates struct {
	LoginPage    *template.Template
	LogoutPage   *template.Template
	RegisterPage *template.Template
}

func ParseTemplates() *HtmlTemplates {
	loginPage := template.Must(template.ParseFS(templateFiles, "login_page.html"))
	logoutPage := template.Must(template.ParseFS(templateFiles, "logout_page.html"))
	registerPage := template.Must(template.ParseFS(templateFiles, "register_page.html"))
	return &HtmlTemplates{loginPage, logoutPage, registerPage}
}
