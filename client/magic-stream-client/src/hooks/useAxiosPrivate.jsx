import axios from "axios";
//import useAuth from "./useAuth";


const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () =>{
  const axiosAuth = axios.create({
  baseURL:apiUrl,
  withCredentials:true,
});
//const {auth,setAuth} = useAuth;

return axiosAuth;

}
export default useAxiosPrivate;

