package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"
)

type GetRegisterPageRenderer interface {
	RenderRegisterPage() []byte
}

type GetRegisterPageHandler struct {
	renderer GetRegisterPageRenderer
}

func NewGetRegisterPageHandler(renderer GetRegisterPageRenderer) *GetRegisterPageHandler {
	return &GetRegisterPageHandler{renderer}
}

func (handler *GetRegisterPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	html := handler.renderer.RenderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

type RegisterPage struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

type HtmxGetRegisterPageView struct {
	registerPage *template.Template
}

func NewHtmxGetRegisterPageView(addUserPage *template.Template) *HtmxGetRegisterPageView {
	return &HtmxGetRegisterPageView{addUserPage}
}

func (view *HtmxGetRegisterPageView) RenderRegisterPage() []byte {
	templateData := RegisterPage{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := view.registerPage.Execute(buffer, templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
