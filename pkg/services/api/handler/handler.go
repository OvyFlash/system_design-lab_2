package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/api/service"
	"net/http"
	"net/url"
)

type Handler struct {
	*settings.Essential
	service *service.Service
}

func NewHandler(e *settings.Essential, service *service.Service) *Handler {
	return &Handler{
		Essential: e,
		service:   service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.simpleProxy(w, r)
}

func (h *Handler) simpleProxy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	s, err := h.service.GetServiceWithEndpoint(path)
	if err != nil {
		h.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	proxyURL := &url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("localhost:%d", s.Port),
		Path:     path,
		RawQuery: r.URL.RawQuery,
	}
	proxyReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		h.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	// CORS headers
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	// Headers
	for k, vv := range r.Header {
		proxyReq.Header.Set(k, vv[0])
	}
	// Make a request
	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		h.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	// Response headers
	for k, vv := range res.Header {
		for _, v := range vv {
			w.Header().Set(k, v)
		}
	}
	// Response status code and body
	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		h.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	// Write response and status code
	w.WriteHeader(res.StatusCode)
	w.Write(body)
}

func (h *Handler) WriteError(w http.ResponseWriter, statusCode int, e string) {
	h.Error().Msg(e)
	bytes, err := json.Marshal(e)
	if err != nil {
		h.Error().Msgf("writing HTTP response: %s", err.Error())
	}
	h.WriteResponse(w, statusCode, "application/json", bytes)
}

func (h *Handler) WriteResponse(w http.ResponseWriter, statusCode int, contentType string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		h.Error().Msgf("writing HTTP response: %s", err.Error())
	}
}
