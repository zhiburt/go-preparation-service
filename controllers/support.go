package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type ApiError struct {
	Err        error
	StatusCode int
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(`{ "error" : "%v" }`, e.Err)
}

func getJson(i interface{}) ([]byte, error) {
	if i == nil {
		return nil, nil
	}
	b, err := json.MarshalIndent(i, " ", "")
	if err != nil {
		return nil, fmt.Errorf("getJson error: %v", err)
	}

	return b, nil
}

func GetOkJSON(i interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	m["data"] = i
	b, err := getJson(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func WriteOkStatus(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusOK)
	(*w).Write([]byte(`{ "status" : "OK" }`))
}

func decode(i interface{}, src map[string][]string) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(i, src)
	if err != nil {
		e := &ApiError{fmt.Errorf("InvalidRequest -> %v", err), http.StatusInternalServerError}
		log.Warning("decoder error: ", e)
		return e
	}
	return nil
}

//valid : responce JSON error
func valid(i interface{}) (bool, error) {
	_, err := govalidator.ValidateStruct(i)
	if err != nil {
		e := &ApiError{fmt.Errorf("InvalidRequest -> %v", err), http.StatusInternalServerError}
		// if allErrors, ok := err.(govalidator.Errors); ok {
		// 	buf := &bytes.Buffer{}
		// 	for _, fld := range allErrors {
		// 		data := []byte(fmt.Sprintf("field %#v\n\n", fld))
		// 		buf.Write(data)
		// 	}
		// 	e = &ApiError{fmt.Errorf("erorr %v, by %v", e, buf.String()), e.StatusCode}
		// }
		log.Warning("govalidator mistake")
		return false, e
	}

	return true, nil
}

func __error_handle(w *http.ResponseWriter, err error) bool {
	defer log.Info("__error_handle call")
	if err != nil {
		if err, ok := err.(*ApiError); ok {
			(*w).WriteHeader(err.StatusCode)
			(*w).Write([]byte(err.Error()))
		} else {
			(*w).WriteHeader(http.StatusInternalServerError)
			(*w).Write([]byte(err.Error()))
		}
		return true
	}

	return false
}

func init() {
	logSettings := logrus.WithFields(logrus.Fields{
		"conrollers": "",
	})
	logSettings.Logger.Level = logrus.InfoLevel
	logSettings.Logger.SetFormatter(&logrus.JSONFormatter{})
	lf, err := os.OpenFile("./logs/_controllers/preparation_log.json.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err == nil {
		logSettings.Logger.SetOutput(lf)
	}
	log = logSettings.Logger

	// //defer lf.Close()
	// log.Logger.SetOutput(lf)
}
