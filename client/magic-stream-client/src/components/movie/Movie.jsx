import { Button } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import useAuth from "../../hooks/useAuth";

const Movie = ({ movie, updateMovieReview }) => {
  const { auth } = useAuth();
  const navigate = useNavigate();

  const handleReviewClick = () => {
    // Navigate to review page for admins
    if (auth && auth.role === "ADMIN") {
      navigate(`/review/${movie.imdb_id}`);
    }
  };

  return (
    <div className="col-md-4 mb-4">
      <div className="card h-100 shadow-sm">
        <div style={{ position: "relative" }}>
          <img
            src={movie.poster_path}
            alt={movie.title}
            className="card-img-top"
            style={{
              objectFit: "contain",
              height: "250px",
              width: "100%",
            }}
          />
        </div>
        <div className="card-body d-flex flex-column">
          <h5 className="card-title">{movie.title}</h5>
          <p className="card-text mb-2">{movie.imdb_id}</p>

          {/* Show Review Button for Admins */}
          {auth && auth.role === "ADMIN" && (
            <Button variant="info" onClick={handleReviewClick}>
              Review
            </Button>
          )}
        </div>

        {/* Show ranking badge if available */}
        {movie.ranking?.ranking_name && (
          <span
            className="badge bg-dark m-3 p-2"
            style={{ fontSize: "1rem" }}
          >
            {movie.ranking.ranking_name}
          </span>
        )}
      </div>
    </div>
  );
};

export default Movie;
