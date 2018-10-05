package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/asaskevich/govalidator"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

//AllPreparationsController Get all preparation
//json responce
func AllPreparationsController(w http.ResponseWriter, r *http.Request) {
	ctxLoger := log.WithFields(log.Fields{
		"hello": "world",
		"other": 123324,
	})
	log.SetFormatter(&log.TextFormatter{})
	ctxLoger.Warningln("asdkalsdlkasd 123132 zxczc", "2313as", time.Now().UTC())
}

type (
	SearchRequest struct {
		Name        string `valid:", optional" schema:"name"`                                             //
		Description string `valid:", optional" schema:"description"`                                      //
		Type        string `valid:"in(cold|psychological|cardiovascular|others), optional" schema:"type"` //
		MaxPrice    string `valid:"numeric, optional" schema:"max"`                                       //
	}

	SearchByTypeRequest struct {
		Type     string `valid:"in(cold|psychological|cardiovascular|others), required" schema:"type"` //
		MaxPrice string `valid:"numeric, optional" schema:"max"`                                       //
	}
)

//PreparationController Get only one preparation
//json responce
func PreparationController(w http.ResponseWriter, r *http.Request) {
	sr := &SearchRequest{}
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
