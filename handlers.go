package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"regexp"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Jenkin Status API !\n")
}

func BuildIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	var builds = RepoShowAllBuilds()
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		panic(err)
	}
}

func RepositoryIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	var builds = RepoShowAllRepos()
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		panic(err)
	}
}

func BuildShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var BuildId int
	var err error
	if BuildId, err = strconv.Atoi(vars["BuildId"]); err != nil {
		panic(err)
	}
	Build := RepoFindBuild(BuildId)
	if Build.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		allowCORS(w)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Build); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

func BuildCreate(w http.ResponseWriter, r *http.Request) {
	var Build Build

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &Build); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		allowCORS(w)
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	re := regexp.MustCompile("([^/]+)\\.git$")
	Build.RepositoryName = re.FindString(Build.RepositoryUrl)

	t := RepoCreateBuild(Build)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func allowCORS(w http.ResponseWriter){
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
}
