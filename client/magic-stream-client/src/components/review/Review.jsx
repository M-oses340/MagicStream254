import { Form, Button, Card } from 'react-bootstrap';
import { useRef, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import useAxiosPrivate from '../../hooks/useAxiosPrivate';
import useAuth from '../../hooks/useAuth';
import Movie from '../movie/Movie';
import Spinner from '../spinner/Spinner';

const Review = () => {
  const [movie, setMovie] = useState({});
  const [loading, setLoading] = useState(false);
  const revText = useRef();
  const { imdb_id } = useParams();
  const { auth } = useAuth();
  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchMovie = async () => {
      setLoading(true);
      try {
        const response = await axiosPrivate.get(`/movie/${imdb_id}`);
        setMovie(response.data);
      } catch (error) {
        console.error('Error fetching movie:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMovie();
  }, []); // <- single GET request

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await axiosPrivate.patch(`/updatereview/${imdb_id}`, {
        admin_review: revText.current.value
      });

      setMovie((prev) => ({
        ...prev,
        admin_review: response.data?.admin_review ?? prev.admin_review,
        ranking: {
          ranking_name: response.data?.ranking_name ?? prev.ranking?.ranking_name
        }
      }));
    } catch (err) {
      console.error('Error updating review:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <Spinner />;

  return (
    <div className="container py-5">
      <h2 className="text-center mb-4">Admin Review</h2>
      <div className="row g-4">
        <div className="col-12 col-md-6">
          <Card className="shadow-sm h-100 p-3">
            <Movie movie={movie} />
          </Card>
        </div>

        <div className="col-12 col-md-6">
          <Card className="shadow-sm h-100 p-4 d-flex flex-column">
            {auth?.role === 'ADMIN' ? (
              <Form onSubmit={handleSubmit} className="d-flex flex-column h-100">
                <Form.Group className="mb-3 flex-grow-1">
                  <Form.Label>Admin Review</Form.Label>
                  <Form.Control
                    ref={revText}
                    as="textarea"
                    rows={10}
                    defaultValue={movie?.admin_review || ''}
                    placeholder="Write your review here..."
                    style={{ resize: 'vertical' }}
                    required
                  />
                </Form.Group>
                <div className="d-flex justify-content-end mt-auto">
                  <Button variant="info" type="submit">
                    Submit Review
                  </Button>
                </div>
              </Form>
            ) : (
              <div className="alert alert-info flex-grow-1">
                {movie?.admin_review || 'No review available'}
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default Review;
