import { useState,useEffect } from "react";
import axiosClient from '../../api/axios.Config';
import Movies from "../movies/movies";

const Home = () =>{
    const [movies,setMovies] = useState([]);
    const [loading,setLoading] = useState(false)
    const [message,setMessage] = useState();

    useEffect(()=>{
        const fetchMovies = async () =>{
            setLoading(true);
            setMessage("")
            try{
                const response = await axiosClient.get('/movies');
                setMovies(response.data);
                if(response.data.length === 0){
                    setMessage('There are currently no movies available ')
                }


            }catch(error){
                console.error('Error fetching movies:',error)
                setMessage("Error fetching message")

            }finally{
                setLoading(false)
            }
            
        }
        fetchMovies();
    }, []);
    return(
        <>
            {loading ? (
                <Spinner/>
            ):  (
                <Movies movies ={movies} updateMovieReview={updateMovieReview} message ={message}/>
            )}
        </>
    );

};
export default Home;
