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
	"strings"

	// MORE ABOUT GCInstances HERE https://github.com/minimum2scp/geco/blob/master/commands.go
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Jenkin Status API !\n")
}

func BuildIndex(w http.ResponseWriter, r *http.Request) {
	var pageSize int = 50
	var err error

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	
	size := r.URL.Query().Get("size")
	if size != ""{
		if pageSize, err = strconv.Atoi(size); err != nil {
			panic(err)
		}
	}

	var builds = RepoShowAllBuilds(pageSize)
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
	Build.RepositoryName = strings.Replace(Build.RepositoryName, ".git", "", -1)

	t := RepoCreateBuild(Build)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	allowCORS(w)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func GCInstances(w http.ResponseWriter, r *http.Request) {
	var err error
	var res []Instance
	var instances []*compute.Instance

	project := "sbtech-pop-poc" // Update Project name

	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		panic(err)
	}
	computeService, err := compute.New(c)
	if err != nil {
		panic(err)
	}
	aggregatedListCall := computeService.Instances.AggregatedList(project)
	for {
		res, err := aggregatedListCall.Do()
		if err != nil {
			panic(err)
			return
		}
		for _, instancesScopedList := range res.Items {
			instances = append(instances, instancesScopedList.Instances...)
		}
		if res.NextPageToken != "" {
			fmt.Fprint(w, "loading more instances with nextPageToken in %s ...", project)
			aggregatedListCall.PageToken(res.NextPageToken)
		} else {
			break
		}
	}
	for _, ins := range instances {
		zone := (func(a []string) string { return a[len(a)-1] })(strings.Split(ins.Zone, "/"))
		machineType := (func(a []string) string { return a[len(a)-1] })(strings.Split(ins.MachineType, "/"))
		internalIP := ins.NetworkInterfaces[0].NetworkIP
		externalIP := ins.NetworkInterfaces[0].AccessConfigs[0].NatIP
		ins_id := strings.Split(ins.Name, "-")[0]
		res = append(res, Instance{ID: ins_id, NAME: ins.Name, ZONE: zone, MACHINE_TYPE: machineType, INTERNAL_IP: internalIP, EXTERNAL_IP: externalIP, STATUS: ins.Status})
	}
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}


func allowCORS(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, X-XSRF-TOKEN")
}
