package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"BuildIndex",
		"GET",
		"/builds",
		BuildIndex,
	},
	Route{
		"BuildCreate",
		"POST",
		"/build",
		BuildCreate,
	},
	Route{
		"BuildShow",
		"GET",
		"/build/{id}",
		BuildShow,
	},
	Route{
		"RepositoryIndex",
		"GET",
		"/repositories",
		RepositoryIndex,
	},
}
