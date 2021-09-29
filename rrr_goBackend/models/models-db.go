package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one movie and err.
func (m *DBModel) Get(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, description, year, release_date, runtime, rating, mpaa_rating, 
	created_at, updated_at, coalesce(poster, '') from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie Movie
	err := row.Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Poster,
	)
	if err != nil {
		return nil, err
	}

	m.PopulateMovieGenres(ctx, &movie, id)

	// query = `select
	// 		mg.id, mg.movie_id, mg.genre_id, g.genre_name
	// 	from
	// 		movies_genres mg
	// 		left join genres g on (g.id = mg.genre_id)
	// 	where
	// 		mg.movie_id = $1

	// `

	// rows, err := m.DB.QueryContext(ctx, query, id)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// genres := make(map[int]string)
	// for rows.Next() {
	// 	var mg MovieGenre
	// 	err := rows.Scan(
	// 		&mg.ID,
	// 		&mg.MovieID,
	// 		&mg.GenreID,
	// 		&mg.Genre.GenreName,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	genres[mg.ID] = mg.Genre.GenreName
	// }

	// movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and err, if any.
func (m *DBModel) All(genre ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, runtime, rating, mpaa_rating, 
	created_at, updated_at, coalesce(poster, '') from movies %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.Poster,
		)
		if err != nil {
			return nil, err
		}
		m.PopulateMovieGenres(ctx, &movie, movie.Id)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *DBModel) PopulateMovieGenres(ctx context.Context, movie *Movie, id int) error {
	query := `select 
			mg.id, mg.movie_id, mg.genre_id, g.genre_name
		from
			movies_genres mg
			left join genres g on (g.id = mg.genre_id)
		where
			mg.movie_id = $1
	
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.Id,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			return err
		}
		genres[mg.Id] = mg.Genre.GenreName
	}
	// rows.Close()

	movie.MovieGenre = genres
	return nil
}

func (m *DBModel) GenresAll() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	genresQuery := `select id, genre_name, created_at, updated_at from genres order by genre_name`

	rows, err := m.DB.QueryContext(ctx, genresQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println(rows)
	defer rows.Close()

	var allGenres []*Genre
	for rows.Next() {
		// genre := new(Genre)
		var genre Genre
		err = rows.Scan(
			&genre.Id,
			&genre.GenreName,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		allGenres = append(allGenres, &genre)
	}
	return allGenres, nil
}

func (m *DBModel) InsertMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	insertQuery := `insert into movies (title, description, year, release_date, runtime, rating,
		mpaa_rating, created_at, updated_at, poster) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := m.DB.ExecContext(ctx, insertQuery,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Poster,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) UpdateMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateQuery := `update movies set title = $1, description = $2, year = $3, release_date = $4, 
		runtime = $5, rating =$6, mpaa_rating = $7, updated_at = $8, poster = $9
		where id = $10`

	_, err := m.DB.ExecContext(ctx, updateQuery,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.CreatedAt,
		movie.Poster,
		movie.Id,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	deleteQuery := `delete from movies where id = $1`

	_, err := m.DB.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
