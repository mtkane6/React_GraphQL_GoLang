package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"rrr_goBackend/models"
	"strings"

	"github.com/graphql-go/graphql"
)

// holds the data we will get from our query
var movies []*models.Movie

// **SCHEMA **
// describe the fields to be used in graphql
var fields = graphql.Fields{ // map[string]*Field

	// "this will be a 'movie', read into this graphql.field{}"   Get me One movie
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{ // map[string]*ArgumentConfig
			"id": &graphql.ArgumentConfig{ // ** This serves to compare if the movie in the DB matches the incoming request
				Type:         graphql.Int, // id is an Int
				DefaultValue: "",
				Description:  "",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int) // check if Args has "id" and cast (type assertion) to int
			if ok {
				for _, movie := range movies {
					if movie.Id == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	// Get me all the movies
	"list": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get all movies",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Search movies by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
				Description:  "",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var theList []*models.Movie
			search, ok := p.Args["titleContains"].(string)
			if ok {
				for _, current := range movies {
					if strings.Contains(current.Title, search) {
						log.Println(fmt.Sprintf("Found match: Search- %s, Movie- %s", search, current.Title))
						theList = append(theList, current)
					}
				}
			}
			return theList, nil
		},
	},
}

var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"poster": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// moviesGraphQL is the handle to get all the movies in the DB
func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling Movies GraphQL")
	movies, _ = app.models.DB.All() // get all the movies

	// get the request from the front end
	q, _ := io.ReadAll(r.Body)

	// query is the actual request statement from the frontend
	query := string(q)
	log.Println("Incoming query: " + query)

	// gets The GraphQL type system to use when validating and executing a query.
	schema, err := getGQLschema()
	if err != nil {
		app.errorJSON(w, errors.New("failed to create schema"))
		log.Println(err)
		return
	}

	// what params are we getting in our request? (what params is the front end asking for?)
	params := graphql.Params{
		Schema:        schema, // The GraphQL type system to use to fulfill request
		RequestString: query,  // the request being made
	}

	// this handles the actual graphql request, returns the query result
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		app.errorJSON(w, fmt.Errorf("failed: %+v", resp.Errors))
	}

	// finally return result to frontend
	j, _ := json.MarshalIndent(resp, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func getGQLschema() (graphql.Schema, error) {
	// specify a root query, states what all the available Fields are for a front end to request
	rootquery := graphql.ObjectConfig{
		Name:   "RootQUery",
		Fields: fields,
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootquery),
	}

	// schema = has a schemaConfig, that has the rootQuery, which has the available fields
	return graphql.NewSchema(schemaConfig)
}
