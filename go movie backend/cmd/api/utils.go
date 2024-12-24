package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (app *application) WriteJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	cont, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(cont)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) ReadJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	max := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(max))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("must contain one json value")
	}
	return nil
}

func (app *application) ErrorJson (w http.ResponseWriter, err error,status ...int) error{
	statusCode:=http.StatusBadRequest
	if len(status)>0{
		statusCode=status[0]
	}
	var payload JSONResponse
	payload.Error=true
	payload.Message=err.Error()
return app.WriteJson(w,statusCode,payload)

}
