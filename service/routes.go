package service

import "net/http"

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the routes
type Routes []Route

var routes = Routes{
	Route{
		"HealthCheck",
		"GET",
		"/limitmonitoring/gethealth",
		getHealth,
	},
}
