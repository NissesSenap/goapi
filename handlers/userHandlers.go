package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/asdine/storm"

	"gopkg.in/mgo.v2/bson"

	"github.com/NissesSenap/GoAPI/cache"
	"github.com/NissesSenap/GoAPI/user"
)

func bodyToUser(r *http.Request, u *user.User) error {
	if r == nil {
		return errors.New("a request is requiered")
	}
	if r.Body == nil {
		return errors.New("request bodi is empty")
	}
	if u == nil {
		return errors.New("a user is requeried")
	}
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bd, u)
}

func usersGetAll(w http.ResponseWriter, r *http.Request) {
	if cache.Serve(w, r) {
		return
	}
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": users})
}

func usersGetOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	if cache.Serve(w, r) {
		return
	}
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})
}

func usersDeleteOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(r))
	w.WriteHeader(http.StatusOK)
}

func usersPutOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	// TODO catch if id exist or not
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(r))
	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})
}

func usersPatchOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}

	err = bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	// u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(r))
	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})
}

func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)
}
