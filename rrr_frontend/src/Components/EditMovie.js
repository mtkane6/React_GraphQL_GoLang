import React, { Component, Fragment } from 'react';
import Input from "./form-components/Input";
import TextInputArea from "./form-components/TextInputArea";
import Select from "./form-components/Select";
import Alert from "./UI-Components/Alert";
import { Link } from 'react-router-dom';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';

export default class EditMovie extends Component {

    constructor(props) {
        super(props);
        this.state = {
            movie: {
                id: 0,
                title: "",
                release_date: "",
                runtime: "",
                mpaa_rating: "",
                description: "",
                rating: "",
            },
            mpaaOptions: [
                {id: "G", value: "G"},
                {id: "PG", value: "PG"},
                {id: "PG-13", value: "PG-13"},
                {id: "R", value: "R"},
                {id: "NC-17", value: "NC-17"},
            ],
            isLoaded: false,
            error: null,
            errors: [],
            alert: {
                type: "d-none",
                message: "",
            },
        }
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit = (evt) => {
        evt.preventDefault();

        // client side form validation.
        let formErrors = [];
        if (this.state.movie.title === "") {
            formErrors.push("title");
        }
        if (this.state.movie.release_date === "") {
            formErrors.push("release_date");
        }
        if (this.state.movie.runtime === "") {
            formErrors.push("runtime");
        }
        if (this.state.movie.rating === "") {
            formErrors.push("rating");
        }
        // if (this.state.movie.description === "") { // text area empty behaves differently for some reason.
        //     formErrors.push("description");
        // }
        this.setState({errors: formErrors})

        if (this.state.errors.length > 0) {
            return false;
        }

        const data = new FormData(evt.target);
        const payload = Object.fromEntries(data.entries()); // gets all the form data after clicking the submit button.
        // console.log(payload);
        const reqOptions = {
            method: 'POST',
            body: JSON.stringify(payload),
        }
        
        fetch("http://localhost:4000/v1/admin/editmovie", reqOptions)
        .then(response => response.json()) // creates jsonData
        .then(jsonData => {
            if (jsonData.error) {
                this.setState({alert: {type:"alert-danger", message: jsonData.error.message}});
            } else {
                this.props.history.push("/admin") // redirect back to "/admin"
            }
        });
    }

    // takes the name and value from the given form input that is handed to this fxn.
    handleChange = (evt) => {
        let value = evt.target.value;
        let name = evt.target.name;

        this.setState((prevState) => ({
            movie: {
                ...prevState.movie,
                [name]: value, // name of current input we are dealing with
            }
        }))
    }

    hasError(key) {
        return this.state.errors.indexOf(key) !== -1;
    }

    componentDidMount () {
        const id = this.props.match.params.id; // once this page is loaded, the const. has been called, get the param from the url.
        console.log("Edit movie id:" + id)
        if (id > 0) {
            fetch("http://localhost:4000/v1/movie/" + id)
            .then(response => {
                if (response.status !== "200" ) {
                    let err = Error;
                    err.Message = "Invalid response code: " +response.status;
                    this.setState({
                        error: err
                    });
                }
                return response.json();
            })
            .then((json) => {
                const releaseDate = new Date(json.movie.release_date);
                this.setState({
                    movie: {
                        id: id,
                        title: json.movie.title,
                        release_date: releaseDate.toISOString().split("T")[0],
                        runtime: json.movie.runtime,
                        mpaa_rating: json.movie.mpaa_rating,
                        rating: json.movie.rating,
                        description: json.movie.description,
                    },
                    isLoaded: true,
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error,
                    })
                })
            })
        } else {
            this.setState({ isLoaded: true});
        }
    
    }

    confirmDelete = (evt) => {
        confirmAlert({
            title: 'Delete Movie',
            message: 'Are you sure?',
            buttons: [
              {
                label: 'Yes',
                onClick: () => {
                    fetch("http://localhost:4000/v1/admin/deletemovie/" + this.state.movie.id, {method: "GET"})
                    .then((response) => response.json())
                    .then(jsonData => {
                        if (jsonData.error) {
                            this.setState({
                                alert: {type: "alert-danger", message: jsonData.error.message}
                            })
                        } else {
                            this.props.history.push({
                                pathname: "/admin",
                            })
                        }
                    });
                }
              },
              {
                label: 'No',
                onClick: () => {}
              }
            ]
          });
    }

    render() {
        let {movie, isLoaded, error} = this.state;
        if (error) {
            return <div>Error: {error.message}</div>
        } else if (!isLoaded) {
            return <p>Loading...</p>
        } else {
            return(
                <Fragment>
                    <h2>Add/Edit Movie</h2>
                    <Alert
                        alertType={this.state.alert.type}
                        alertMessage={this.state.alert.message}
                    />
                    <hr/>
                    <form onSubmit={this.handleSubmit}>
                        <input 
                            type = "hidden"
                            name="id"
                            id="id"
                            value={movie.id}
                            onChange={this.handleChange}
                        />
                        {/* <div className="mb-3">
                            <label htmlFor="title" className="form-label">
                                Title
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                id="title"
                                name="title"
                                value={movie.title}
                                onChange={this.handleChange}
                            />
                        </div> */}
                        <Input
                            title={"Title"}
                            className={this.hasError('title') ? "is-invalid" : ""} // if hasError == true, styles text box as error
                            type={'text'}
                            name={'title'}
                            value={movie.title}
                            handleChange={this.handleChange}
                            errorDiv={this.hasError('title') ? "text-danger" : "d-none"} // for styling error message
                            errorMsg={"Please enter a Title"}
                        />
                        {/* <div className="mb-3">
                            <label htmlFor="release_date" className="form-label"> {/*label is for screen readers, accessability features*/}{/*}
                                Release Date
                            </label>
                            <input
                                type="text"
                                className="form-control" // for css
                                id="release_date" // for css
                                name="release_date"  // this matches with handleChange().  '[name]'. Must match spelling in this.state.
                                value={movie.release_date} // what is displayed in the text box.
                                onChange={this.handleChange} // takes value from box, and updates the matching state value.
                            />
                        </div> */}
                        <Input
                            title={"Release Date"}
                            className={this.hasError('release_date') ? "is-invalid" : ""}
                            type={'date'}
                            name={'release_date'}
                            value={movie.release_date}
                            handleChange={this.handleChange}
                            errorDiv={this.hasError('release_date') ? "text-danger" : "d-none"}
                            errorMsg={"Please enter a Date"}
                        />
                        {/* <div className="mb-3">
                            <label htmlFor="runtime" className="form-label">
                                Run Time
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                id="runtime"
                                name="runtime"
                                value={movie.runtime}
                                onChange={this.handleChange} 
                            />
                        </div> */}
                        <Input
                            title={"Run Time"}
                            className={this.hasError('runtime') ? "is-invalid" : ""}
                            type={'text'}
                            name={'runtime'}
                            value={movie.runtime}
                            handleChange={this.handleChange}
                            errorDiv={this.hasError('runtime') ? "text-danger" : "d-none"} // for styling error message
                            errorMsg={"Please enter a Runtime"}
                        />
                        
                        {/* <div className="mb-3">
                            <label htmlFor="mpaa_rating" className="form-label">
                                MPAA Rating
                            </label>
                            <select name="mpaa_rating"  /* <- this matches with handleChange().  '[name]'. Must match spelling in this.state.*/ /*
                                className="form-select" 
                                value={movie.mpaa_rating} 
                                onChange={this.handleChange}>

                                <option className="form-select">Choose...</option>
                                <option className="form-select" value={"G"}>G</option>
                                <option className="form-select" value={"PG"}>PG</option>
                                <option className="form-select" value={"PG-13"}>PG-13</option>
                                <option className="form-select" value={"R"}>R</option>
                                <option className="form-select" value={"NC-17"}>NC-17</option>
                            </select>
                        </div> */}

                        <Select
                            title={'MPAA Rating'}
                            name={'mpaa_rating'}
                            options={this.state.mpaaOptions}
                            values={movie.mpaa_rating}
                            handleChange={this.handleChange}
                            placeholder={'Choose...'}
                        />

                        <Input
                            title={"Rating"}
                            className={this.hasError('rating') ? "is-invalid" : ""}
                            type={'text'}
                            name={'rating'}
                            value={movie.rating}
                            handleChange={this.handleChange}
                            errorDiv={this.hasError('rating') ? "text-danger" : "d-none"} // for styling error message
                            errorMsg={"Please enter a Rating"}
                        />

                        {/* <div className="mb-3">
                            <label htmlFor="rating" className="form-label">
                                Rating
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                id="rating"
                                name="rating"
                                value={movie.rating}
                                onChange={this.handleChange} 
                            />
                        </div> */}

                        {/* <div className="mb-3">
                            <label htmlFor="description" className="form-label">
                                Description
                            </label>
                            <textarea 
                                className="form-control" 
                                id="description" 
                                name="description"
                                rows="3" 
                                value={movie.description}
                                onChange={this.handleChange}
                            />
                        </div> */}
                        <TextInputArea // these 'Dumb' components must be self-closing tags. <tag/>, not <tag></tag>
                            title={"Description"}
                            className={this.hasError('description') ? "is-invalid" : ""}
                            name={"description"}
                            rows={"3"}
                            value={movie.description}
                            handleChange={this.handleChange}
                            errorDiv={this.hasError('description') ? "text-danger" : "d-none"} // for styling error message
                            errorMsg={"Please enter a Description"}
                        />

                        <hr/>
                        <button className="btn btn-primary">
                            Save
                        </button>
                        <Link to="/admin" className="btn btn-warning ms-1">
                            Cancel
                        </Link>
                        {movie.id > 0 && ( /*inline conditional expression: if TRUE && func => func.  If FALSE && func => FALSE */
                            <a href="#!" onClick={() => this.confirmDelete()} className="btn btn-danger ms-1">
                                Delete
                            </a>
                        )}
                    </form>

{/* // This shows live values of the component state */}
                    {/* <div className="mt-3">
                        <pre>{JSON.stringify(this.state, null, 3)}</pre>
                    </div> */}
                </Fragment>
            );
        }
    }
}