package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"
)

type GetLoginPageRenderer interface {
	RenderLoginPage() []byte
}

type GetLoginPageHandler struct {
	renderer GetLoginPageRenderer
}

func NewGetLoginPageHandler(renderer GetLoginPageRenderer) *GetLoginPageHandler {
	return &GetLoginPageHandler{renderer}
}

func (handler *GetLoginPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	html := handler.renderer.RenderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

type LoginPage struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

type HtmxGetLoginPageView struct {
	loginPage *template.Template
}

func NewHtmxGetLoginPageView(loginPage *template.Template) *HtmxGetLoginPageView {
	return &HtmxGetLoginPageView{loginPage}
}

func (view *HtmxGetLoginPageView) RenderLoginPage() []byte {
	templateData := LoginPage{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := view.loginPage.Execute(buffer, templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
