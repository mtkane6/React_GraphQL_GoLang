import React, { Component, Fragment } from 'react';
// import {Link} from 'react-router-dom';

export default class OneMovie extends Component {
    state = {movie: {}, isLoaded: false, error: null};

    componentDidMount() {
        fetch("http://localhost:4000/v1/movie/"+this.props.match.params.id)
        // .then((response) => response.json())
        .then((response) => {
            console.log("Status code is:", response.status);
            if (response.status !== 200) {
                let err = Error;
                err.message = "Invalid response code:" + response.status;
                this.setState({error: err});
            }
            return response.json();
        })
        .then((json)=> {
            this.setState({
                movie: json.movie,
                isLoaded: true
            },
            (error) => {
                this.setState({
                    isLoaded: true,
                    error: null,
                });
            }
            );
        });
    }

    render() {
        const {movie, isLoaded, error} = this.state;
        if (movie.genres) {
            movie.genres= Object.values(movie.genres)
        } else {
            movie.genres = [];
        }

        if (error) {
            return <div>Error: {error.message} </div>
        }
        else if(!isLoaded) {
            return <p>Loading...</p>
        }
        else {
        return (
            <Fragment>
                <h2>Movie: {movie.title} ({movie.year})</h2>
                <div className="float-start">
                    <small>Rating: {movie.mpaa_rating}</small>
                </div>
                <div className="float-end">
                    {movie.genres.map((m,index) => (
                        <span className="badge bg-primary me-1" key={index}> {/* key necessary when interating through map*/}
                            {m}
                        </span>
                    ))}
                </div>
                <div className="clearfix"></div>
                <hr/>

                <table className="table table-compact table-striped">
                    <thead></thead>
                    <tbody>
                        <tr>
                            <td>
                                <strong>Title:</strong>
                            </td>
                            <td>{movie.title}</td>
                        </tr>
                        <tr>
                            <td><strong>Description:</strong></td>
                            <td>{movie.description}</td>
                        </tr>
                        <tr>
                            <td>
                                <strong>RunTime:</strong>
                            </td>
                            <td>{movie.runtime}</td>
                        </tr>
                    </tbody>
                </table>
            </Fragment>
        );
        }
    }
}