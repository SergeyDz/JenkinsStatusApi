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
	"google.golang.org/api/compute/v1"
	"time"
	"net"
	"os"
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
	if size != "" {
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
	var res []Instance
	var ign_list []string
	var instances []*compute.Instance

	project := "sbtech-pop-poc" // Update Project name
	ign := strings.Split(os.Getenv("IgnorEnv"), ",")
	for x := range ign {
		ign_list = append(ign_list, ign[x])
	}

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
	jenlist := GCBuildStatus()
	for _, ins := range instances {
		ins_id := strings.Split(ins.Name, "-")[0]
		if stringInSlice(ins_id, ign_list) {
			continue
		}
		zone := (func(a []string) string { return a[len(a)-1] })(strings.Split(ins.Zone, "/"))
		machineType := (func(a []string) string { return a[len(a)-1] })(strings.Split(ins.MachineType, "/"))
		internalIP := ins.NetworkInterfaces[0].NetworkIP
		externalIP := ins.NetworkInterfaces[0].AccessConfigs[0].NatIP

		inst_app := strings.Split(ins.Name, "-")
		jenstat := "false"
		for i := range jenlist {
			jen_el := strings.Split(jenlist[i], ";")
			if (jen_el[0] == ins_id) && (jen_el[1] == inst_app[len(inst_app)-1]) {
				jenstat = jen_el[2]
				break
			}
		}
		ostype := "unknown"
		for _, inmeta := range ins.Metadata.Items {
			ostype = strings.Split(inmeta.Key, "-")[0]
		}
		res = append(res, Instance{ID: ins_id, NAME: ins.Name, ZONE: zone, MACHINE_TYPE: machineType, OS: ostype, INTERNAL_IP: internalIP, EXTERNAL_IP: externalIP, STATUS: ins.Status, JENBUILD: jenstat})
	}
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func GCBuildStatus() []string {
	var jobs JenkinsBuilds
	var list []string

	json.Unmarshal(getHTTPBody("http://jenkins.paas.sbtech.com:8080/job/Common/job/Create_application_terraform_poc_test/api/json?tree=builds[url,building,actions[parameters[name,value]]]"), &jobs)
	for _, val := range jobs.Builds {
		if val.Building {
			for i := range val.Actions[0].Parameters {
				if val.Actions[0].Parameters[i].Name == "application" {
					for y := range val.Actions[0].Parameters {
						if val.Actions[0].Parameters[y].Name == "env_name" {
							list = append(list, (val.Actions[0].Parameters[y].Value + ";" + val.Actions[0].Parameters[i].Value + ";" + val.Url))
							break
						}
					}
				}
			}
		}
	}
	return list
}

func Ping(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostName := vars["url"]
	portNum := vars["port"]
	typ := vars["type"]
	seconds, _ := strconv.Atoi(vars["timeout"])
	timeOut := time.Duration(seconds) * time.Second

	allowCORS(w)
	w.WriteHeader(http.StatusOK)

	start := time.Now()
	if typ == "http" {
		httpClient := http.Client{
			Timeout: timeOut,
		}
		resp, err := httpClient.Get("http://" + hostName + ":" + portNum)
		if (err != nil || resp == nil) {
			if err := json.NewEncoder(w).Encode("0"); err != nil {
				panic(err)
			}
			return
		}
		defer resp.Body.Close()
	} else {
		conn, err := net.DialTimeout("tcp", hostName+":"+portNum, timeOut)
		if err != nil {
			if err := json.NewEncoder(w).Encode("0"); err != nil {
				panic(err)
			}
			return
		}
		conn.Close()
	}
	elapsed := time.Since(start)
	ping := elapsed.Nanoseconds() / int64(time.Millisecond)
	if err := json.NewEncoder(w).Encode(ping); err != nil {
		panic(err)
	}

}

func JenkinsJobs(w http.ResponseWriter, r *http.Request) {
	var appterrpoc AppTerraformPOC
	var buildartpop BuilsArtifacrtPOP
	var cloudenvpoc CloudEnvPOC
	var list []JenkinsJob

	json.Unmarshal(getHTTPBody("http://jenkins.paas.sbtech.com:8080/job/Common/job/Create_application_terraform_poc_test/api/json?tree=builds[displayName,result,url,building,timestamp,actions[causes[upstreamUrl,upstreamBuild]]]"), &appterrpoc)
	for _, val := range appterrpoc.Builds {
		joburl := ""
		jobid := ""
		jobUserName := ""
		for i := range val.Actions {
			if val.Actions[i].Class == "hudson.model.CauseAction" {
				if val.Actions[i].Causes[0].Class == "hudson.model.Cause$UpstreamCause" {
					joburl = val.Actions[i].Causes[0].UpstreamURL
					jobid = strconv.Itoa(val.Actions[i].Causes[0].UpstreamBuild)
					break
				}
			}
		}
		json.Unmarshal(getHTTPBody("http://jenkins.paas.sbtech.com:8080/"+joburl+jobid+"/api/json?tree=actions[causes[upstreamUrl,upstreamBuild]]"), &buildartpop)
		for i := range buildartpop.Actions {
			if buildartpop.Actions[i].Class == "hudson.model.CauseAction" {
				if buildartpop.Actions[i].Causes[0].Class == "hudson.model.Cause$UpstreamCause" {
					joburl = buildartpop.Actions[i].Causes[0].UpstreamURL
					jobid = strconv.Itoa(buildartpop.Actions[i].Causes[0].UpstreamBuild)
					break
				}
			}
		}
		json.Unmarshal(getHTTPBody("http://jenkins.paas.sbtech.com:8080/"+joburl+jobid+"/api/json?tree=actions[causes[userName]]"), &cloudenvpoc)
		for i := range cloudenvpoc.Actions {

			if cloudenvpoc.Actions[i].Class == "hudson.model.CauseAction" {
				if cloudenvpoc.Actions[i].Causes[0].Class == "hudson.model.Cause$UserIdCause" {
					jobUserName = cloudenvpoc.Actions[i].Causes[0].UserName
					break
				}
			}
		}
		list = append(list, JenkinsJob{Timestamp: time.Unix(val.Timestamp/1000, 0).Format("15:04"),
			Building: val.Building,
			Result: val.Result,
			DisplayName: val.DisplayName,
			URL: val.URL,
			UserName: jobUserName})
	}
	allowCORS(w)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func getHTTPBody(url string) []byte {
	var client http.Client

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+os.Getenv("Jenkins64base"))
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	return body
}

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, X-XSRF-TOKEN")
}
