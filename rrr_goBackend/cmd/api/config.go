package main

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

var apiKey string = "554b471f2e604ca475cf4843efffd7c4" // I know this is bad practice, but this is just a demo.

var dbURL string = "https://api.themoviedb.org/3/search/movie?api_key="
