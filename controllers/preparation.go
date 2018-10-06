package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"../db"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type ApiError struct {
	Err        error
	StatusCode int
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(`"error" : "%v"`, e.Err)
}

//AllPreparationsController Get all preparation
//json responce
func AllPreparationsController(w http.ResponseWriter, r *http.Request) {
	allPreparations := db.GetAllPreparations()
	// if allPreparations == nil {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`"data": null`))
	// 	return
	// }

	m := make(map[string][]*db.Preparation)
	m["data"] = allPreparations
	b, err := getOkJSON(allPreparations)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

type (
	FindPreparationByNameRequest struct {
		Name string `valid:", required" schema:"name"`
	}
	FindPreparationByIdRequest struct {
		Id int `valid:"numeric, required" schema:"id"`
	}

	PreparationSearchRequest struct {
		Name             string `valid:", optional" schema:"name"`                                             //
		Description      string `valid:", optional" schema:"description"`                                      //
		Type             string `valid:"in(cold|psychological|cardiovascular|others), optional" schema:"type"` //
		ActiveIngredient string `valid:", optional" schema:"activeIngredient"`
		MaxPrice         int    `valid:"numeric, optional" schema:"max"` //
	}

	SearchByTypeRequest struct {
		Type     string `valid:"in(cold|psychological|cardiovascular|others), required" schema:"type"` //
		MaxPrice string `valid:"numeric, optional" schema:"max"`                                       //
	}
)

func FindPreparationByIdController(w http.ResponseWriter, r *http.Request) {
	sr := &FindPreparationByIdRequest{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	if sr.Id < 0 {
		e := &ApiError{fmt.Errorf("InvalidRequest -> id must be > 0", err), http.StatusInternalServerError}
		w.WriteHeader(e.StatusCode)
		w.Write([]byte(e.Err.Error()))
		return
	}

	preparation := db.FindPreparationById(sr.Id)

	b, err := getOkJSON(preparation)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func FindPreparationByNameController(w http.ResponseWriter, r *http.Request) {
	sr := &FindPreparationByNameRequest{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	if sr.Name == "" {
		e := &ApiError{fmt.Errorf("InvalidRequest : name == \"\"", err), http.StatusInternalServerError}
		w.WriteHeader(e.StatusCode)
		w.Write([]byte(e.Err.Error()))
		return
	}

	preparation := db.FindPreparationByName(sr.Name)

	b, err := getOkJSON(preparation)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//PreparationController Get only one preparation
//json responce
func PreparationController(w http.ResponseWriter, r *http.Request) {
	sr := &PreparationSearchRequest{}
	vars := mux.Vars(r)
	sr.Name = vars["name"]

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(sr, r.URL.Query())
	if err != nil {
		log.Println("decoder error: ", err)
		http.Error(w, "iternal server error :(", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(govalidator.ToString(sr) + "\n"))

	_, err = govalidator.ValidateStruct(sr)
	if err != nil {

		if allErrors, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrors {
				data := []byte(fmt.Sprintf("field %#v\n\n", fld))
				w.Write(data)
			}
		}
		log.Println("govalidator mistake")
		return
	}

	w.Write([]byte("HELLO WORLD"))
}

func PreparationsByTypeController(w http.ResponseWriter, r *http.Request) {
	sr := &SearchByTypeRequest{}
	vars := mux.Vars(r)
	sr.Type = vars["type_name"]

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(sr, r.URL.Query())
	if err != nil {
		log.Println("decoder error: ", err)
		http.Error(w, "iternal server error :(", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(govalidator.ToString(sr) + "\n"))

	_, err = govalidator.ValidateStruct(sr)
	if err != nil {

		if allErrors, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrors {
				data := []byte(fmt.Sprintf("field %#v\n\n", fld))
				w.Write(data)
			}
		}
		log.Println("govalidator mistake")
		return
	}

	w.Write([]byte("HELLO WORLD by type"))
}

func getJson(i interface{}) ([]byte, error) {
	if i == nil {
		return nil, nil
	}
	b, err := json.MarshalIndent(i, "", "")
	if err != nil {
		return nil, fmt.Errorf("getJson error: %v", err)
	}

	return b, nil
}

func getOkJSON(i interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	m["data"] = i
	b, err := getJson(m)
	if err != nil {
		return nil, err
	}

	return b, nil
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
		"conroller": "preparation",
	})
	logSettings.Logger.Level = logrus.InfoLevel
	logSettings.Logger.SetFormatter(&logrus.JSONFormatter{})
	lf, err := os.OpenFile("./logs/_controllers/preparation_log.json.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	logSettings.Logger.SetOutput(lf)
	log = logSettings.Logger

	// //defer lf.Close()
	// log.Logger.SetOutput(lf)
}
