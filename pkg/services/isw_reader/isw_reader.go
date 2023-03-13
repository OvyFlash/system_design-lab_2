package isw_reader

import (
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/isw_reader/service"
)

// get raw data from isw
// write to db
type ISWReader struct {
	service *service.Service

	*settings.Essential
}

func NewISWReader(e *settings.Essential) *ISWReader {
	i := &ISWReader{
		Essential: e,
	}
	i.setHTTPHandlers()
	i.service = service.NewService(e, e.Logger)
	return i
}

func (i *ISWReader) Start() (err error) {
	i.service.Start()
	return nil
}

func (i *ISWReader) Stop() error {
	return nil
}

func (i *ISWReader) setHTTPHandlers() {
	// i.Router.GET("/getWeather", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	w.Write([]byte("hello world"))
	// })
}
