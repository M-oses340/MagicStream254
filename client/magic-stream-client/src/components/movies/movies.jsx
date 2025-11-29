import Movie from "../Movie";

const Movies = ({movies,message}) => {
    return(
        <div className="container mt-4">
            <div className="row">
                {movies && movies.len>0
                  ? movies.map((movie)=>(
                    <Movie key={movie_id} movie={movie}/>
                  ))
                   : <h2>{message}</h2>
                }

            </div>

        </div>
    )
}