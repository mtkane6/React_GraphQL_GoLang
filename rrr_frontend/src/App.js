import React, {/*Fragment*/} from "react";
import {BrowserRouter as Router, Switch, Route, Link, /*useParams, useRouteMatch*/} from 'react-router-dom';
import Movies from './Components/Movies';
import Admin from './Components/Admin';
import Home from './Components/Home';
// import Categories from "./Components/Categories";
import Genres from './Components/Genres';
import OneMovie from "./Components/OneMovie";
import OneGenre from "./Components/OneGenre";
import EditMovie from "./Components/EditMovie";
import GraphQL from "./Components/GraphQL";
import OneMovieGraphQL from "./Components/OneMovieGraphQL";

export default function App() {
  return (
    <Router>
    <div className="container">
      <div className="row">
        <h1 className="mt-3">
          Go Watch A Movie
        </h1>
        <hr className="mb-3"></hr>
      </div>

      <div className="row">
        <div className="col-md-2">
          <nav>
            <ul className="list-group">
              <li className="list-group-item">
                <Link to="/">Home</Link>
              </li>
              <li className="list-group-item">
                <Link to="/movies">Movies</Link>
              </li>
              <li className="list-group-item">
                <Link to="/genres">Genres</Link>
              </li>
              <li className="list-group-item">
                <Link to="/admin/movie/0">Add Movie</Link>
              </li>
              <li className="list-group-item">
                <Link to="/graphql">GraphQL</Link>
              </li>
              <li className="list-group-item">
                <Link to="/admin">Manage Catalogue</Link>
              </li>
            </ul>
          </nav>
          
        </div>
        <div className="col-md-10">
          <Switch>
            {/* <Route path="/movies/:id">
              <Movie />
            </Route> */}
            <Route path="/movies/:id" component={OneMovie} />

            <Route path="/movies">
              <Movies />
            </Route>

            <Route path="/genre/:id" component={OneGenre} />
            <Route exact path="/genres">
              <Genres />
            </Route>

            {/* <Route exact path="/by-category/comedy">
              <Categories title={`Comedy`}/>
            </Route> */}
            {/* <Route exact path="/by-category/comedy"
            render={
              (props) => <Categories {...props} title={`Comedy`}/>
            }
            />

            <Route exact path="/by-category/drama"
            render={
              (props) => <Categories {...props} title={`Drama`}/>
            }
            /> */}
            <Route path="/admin/movie/:id" component={EditMovie}/>
            <Route path="/admin">
              <Admin />
            </Route>
            <Route exact path="/graphql">
              <GraphQL />
            </Route>
            
            <Route path="/moviesgraphql/:id" component={OneMovieGraphQL} />
            <Route path="/">
              <Home />
            </Route>
          </Switch>

        </div>
      </div>
    </div>
    </Router>
  );
}

// function Movie() {
//   let { id } = useParams();
//   return <h2>Movie id: {id}</h2>
// }

// function CategoryPage() {
//   let { path, /*url*/ } = useRouteMatch();
//   return (
//     <div>
//       <h2>
//         Categories
//       </h2>
//       <ul>
//         <li>
//           <Link to={`${path}/comedy`}>Comedy</Link>
//         </li>
//         <li>
//           <Link to={`${path}/drama`}>Drama</Link>
//         </li>
//       </ul>
//     </div>
//   );
// }
