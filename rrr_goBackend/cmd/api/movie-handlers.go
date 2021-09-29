package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"rrr_goBackend/models"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("Calling 'Get One Movie'")
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	// movie := models.Movie{
	// 	ID:          id,
	// 	Title:       "Some movie",
	// 	Description: "Some desciption",
	// 	Year:        2021,
	// 	ReleaseDate: time.Date(2021, 01, 01, 01, 0, 0, 0, time.Local),
	// 	Runtime:     100,
	// 	Rating:      5,
	// 	MPAARating:  "PG-13",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	movie, err := app.models.DB.Get(id)
	if err != nil {
		app.errorJSON(w, err)
	}

	_ = app.writeJSON(w, http.StatusOK, movie, "movie")
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("Calling 'All Movies'")
	movies, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
	}

	_ = app.writeJSON(w, http.StatusOK, movies, "movies")
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("Calling 'All Genres'")
	var genres []*models.Genre
	genres, err := app.models.DB.GenresAll()
	if err != nil {
		app.errorJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("Calling 'Get All Movies By Genre'")
	params := httprouter.ParamsFromContext(r.Context())

	genreId, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movies, err := app.models.DB.All(genreId)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.logger.Printf("Got movies by genre: %+v", movies)
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling Delete Movie")

	params := httprouter.ParamsFromContext(r.Context())
	movieId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.models.DB.DeleteMovie(movieId)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.statusResponse(w, true, "response")
}

type moviePayload struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (app *application) editMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling editMovie")
	var payload moviePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var movie models.Movie

	if payload.Id != "0" {
		id, _ := strconv.Atoi(payload.Id)
		m, _ := app.models.DB.Get(id)
		movie = *m
		// movie.UpdatedAt = time.Now()
	}

	movie.Id, _ = strconv.Atoi(payload.Id)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	movie.UpdatedAt, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	if movie.Poster == "" {
		movie = getPoster(movie)
	}

	if movie.Id == 0 {
		err = app.models.DB.InsertMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	app.statusResponse(w, true, "response")
}

func getPoster(movie models.Movie) models.Movie {

	client := &http.Client{}

	reqString := dbURL + apiKey + "&query=" + url.QueryEscape(movie.Title)

	log.Println(reqString)

	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Content-Type", "Application/json")
	req.Header.Add("Accept", "Application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var respObj models.TheMovieDB
	err = json.Unmarshal(bodyBytes, &respObj)
	if err != nil {
		log.Println(err)
		return movie
	}

	if len(respObj.Results) > 0 {
		movie.Poster = respObj.Results[0].PosterPath
	}

	return movie
}

func (app *application) statusResponse(w http.ResponseWriter, t_f bool, msg string) {
	ok := jsonResponse{
		OK: t_f,
	}

	err := app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
	}
}
