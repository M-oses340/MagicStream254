import { useState, useEffect } from "react";
import axiosClient from '../../api/axios.Config';
import Movies from "../movies/Movies";

const Home = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  // Placeholder function for movie review updates
  const updateMovieReview = (movieId, review) => {
    console.log(`Update review for ${movieId}:`, review);
    // You can implement the actual update logic later
  };

  useEffect(() => {
    const fetchMovies = async () => {
      setLoading(true);
      setMessage("");
      try {
        const response = await axiosClient.get('/movies');
        setMovies(response.data);
        if (response.data.length === 0) {
          setMessage('There are currently no movies available.');
        }
      } catch (error) {
        console.error('Error fetching movies:', error);
        setMessage("Error fetching movies");
      } finally {
        setLoading(false);
      }
    };
    fetchMovies();
  }, []);

  return (
    <div>
      {loading && <p>Loading movies...</p>}
      {!loading && movies.length > 0 && (
        <Movies movies={movies} updateMovieReview={updateMovieReview} message={message} />
      )}
      {!loading && movies.length === 0 && <p>{message}</p>}
    </div>
  );
};

export default Home;
