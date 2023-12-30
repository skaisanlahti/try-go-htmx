package htmx

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

type extraHeaders = map[string]string

type defaultRenderer struct {
	template *template.Template
}

func newDefaultRenderer(t *template.Template) *defaultRenderer {
	return &defaultRenderer{t}
}

func (this *defaultRenderer) render(response http.ResponseWriter, block string, data any, headers extraHeaders) {
	buffer := &bytes.Buffer{}
	if err := this.template.ExecuteTemplate(buffer, block, data); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Add("Content-type", "text/html; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			response.Header().Set(key, value)
		}
	}

	response.Write(buffer.Bytes())
	return
}

func newRenderKey() int64 {
	return time.Now().UnixMilli()
}
