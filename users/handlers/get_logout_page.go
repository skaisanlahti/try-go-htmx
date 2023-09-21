package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type GetLogoutPageRenderer interface {
	RenderLogoutPage() []byte
}

type GetLogoutPageHandler struct {
	renderer GetLogoutPageRenderer
}

func NewGetLogoutPageHandler(renderer GetLogoutPageRenderer) *GetLogoutPageHandler {
	return &GetLogoutPageHandler{renderer}
}

func (handler *GetLogoutPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	html := handler.renderer.RenderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

type HtmxGetLogoutPageView struct {
	logoutPage *template.Template
}

func NewHtmxGetLogoutPageView(logoutPage *template.Template) *HtmxGetLogoutPageView {
	return &HtmxGetLogoutPageView{logoutPage}
}

func (view *HtmxGetLogoutPageView) RenderLogoutPage() []byte {
	buffer := &bytes.Buffer{}
	err := view.logoutPage.Execute(buffer, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}
