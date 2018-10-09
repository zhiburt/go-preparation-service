package controllers

import (
	"fmt"
	"net/http"

	"uni/coorse/db"
	// "../db"
)

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
	b, err := GetOkJSON(allPreparations)
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

	InsertPreparationRequest struct {
		Name             string `valid:", required" schema:"name"`
		Description      string `valid:", optional" schema:"description"`
		Type             string `valid:", optional" schema:"type"`
		ActiveIngredient string `valid:", optional" schema:"activeIngredient"`
		ImageURL         string `valid:", optional" schema:"imageURL"`
	}

	UpdatePreparationRequest struct {
		Id               int    `valid:"numeric, optional" schema:"id"`
		Name             string `valid:", optional" schema:"name"`
		Description      string `valid:", optional" schema:"description"`
		Type             string `valid:", optional" schema:"type"`
		ActiveIngredient string `valid:", optional" schema:"activeIngredient"`
		ImageURL         string `valid:", optional" schema:"imageURL"`
	}

	DeletePreparationRequest struct {
		Id               int    `valid:"numeric, optional" schema:"id"`
		Name             string `valid:", optional" schema:"name"`
		Description      string `valid:", optional" schema:"description"`
		Type             string `valid:", optional" schema:"type"`
		ActiveIngredient string `valid:", optional" schema:"activeIngredient"`
		ImageURL         string `valid:", optional" schema:"imageURL"`
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
	// if sr.Id < 0 {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest -> id must be > 0"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Err.Error()))
	// 	return
	// }

	preparation := db.FindPreparationById(sr.Id)

	b, err := GetOkJSON(preparation)
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
	// if sr.Name == "" {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest : name == \"\""), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Err.Error()))
	// 	return
	// }

	preparation := db.FindPreparationByName(sr.Name)

	b, err := GetOkJSON(preparation)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func InsertPreparationController(w http.ResponseWriter, r *http.Request) {
	sr := &InsertPreparationRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Name == "" {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest : name must be not empty"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Err.Error()))
	// 	return
	// }

	err = db.InsertIntoPreparations(&db.Preparation{
		Name:             sr.Name,
		Type:             sr.Type,
		ActiveIngredient: sr.ActiveIngredient,
		ImageURL:         sr.ImageURL,
		Description:      sr.Description,
	})

	if err != nil {
		err := &ApiError{fmt.Errorf("error -> %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : "OK" }`))
}

func UpdatePreparationController(w http.ResponseWriter, r *http.Request) {
	sr := &UpdatePreparationRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//by name or id
	if sr.Name == "" && sr.Id == 0 {
		e := &ApiError{fmt.Errorf("InvalidRequest -> name or id must be not empty"), http.StatusInternalServerError}
		w.WriteHeader(e.StatusCode)
		w.Write([]byte(e.Error()))
		return
	}

	stat := db.UpdatePreparation(&db.Preparation{
		Id:               sr.Id,
		Name:             sr.Name,
		Type:             sr.Type,
		ActiveIngredient: sr.ActiveIngredient,
		ImageURL:         sr.ImageURL,
		Description:      sr.Description,
	})

	if !stat {
		err := &ApiError{fmt.Errorf("preparation with name %s or id %d doesn't exists in db", sr.Name, sr.Id), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : "OK" }`))
}

func DeletePreparationController(w http.ResponseWriter, r *http.Request) {
	sr := &DeletePreparationRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//by name or id
	if sr.Name == "" && sr.Id == 0 {
		e := &ApiError{fmt.Errorf("InvalidRequest -> name or id must be not empty"), http.StatusInternalServerError}
		w.WriteHeader(e.StatusCode)
		w.Write([]byte(e.Error()))
		return
	}

	changes := db.DeletePreparation(&db.Preparation{
		Id:               sr.Id,
		Name:             sr.Name,
		Type:             sr.Type,
		ActiveIngredient: sr.ActiveIngredient,
		ImageURL:         sr.ImageURL,
		Description:      sr.Description,
	})

	if changes == 0 {
		err := &ApiError{fmt.Errorf("preparation with name %s or id %d doesn't exists in db", sr.Name, sr.Id), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : "OK" }`))
}
