import { useEffect } from 'react';
import axios from 'axios';
import useAuth from './useAuth';

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {

    const axiosAuth = axios.create({
        baseURL: apiUrl,
        withCredentials: true, // important for cookies
    });

    const { auth, setAuth } = useAuth();

    let isRefreshing = false;
    let failedQueue = [];

    // ---- Helper for queued requests ----
    const processQueue = (error, token = null) => {
        failedQueue.forEach(prom => {
            if (error) {
                prom.reject(error);
            } else {
                prom.resolve(token);
            }
        });
        failedQueue = [];
    };

    useEffect(() => {

        axiosAuth.interceptors.request.use(
            (config) => {
                if (auth?.accessToken) {
                    config.headers["Authorization"] = `Bearer ${auth.accessToken}`;
                }
                return config;
            },
            (error) => Promise.reject(error)
        );

        axiosAuth.interceptors.response.use(
            (response) => response,

            async (error) => {
                console.log("⚠ Interceptor caught:", error);

                const originalRequest = error.config;

                // If refresh token is invalid → logout
                if (
                    originalRequest.url.includes('/refresh') &&
                    error?.response?.status === 401
                ) {
                    console.error("❌ Refresh token expired/invalid");
                    setAuth(null);
                    localStorage.removeItem("user");
                    return Promise.reject(error);
                }

                // --- Handle 401 unauthorized (access token expired) ---
                if (error?.response?.status === 401 && !originalRequest._retry) {

                    if (isRefreshing) {
                        return new Promise((resolve, reject) => {
                            failedQueue.push({ resolve, reject });
                        })
                        .then((token) => {
                            originalRequest.headers["Authorization"] = `Bearer ${token}`;
                            return axiosAuth(originalRequest);
                        })
                        .catch(err => Promise.reject(err));
                    }

                    originalRequest._retry = true;
                    isRefreshing = true;

                    try {
                        // Call refresh endpoint
                        const res = await axiosAuth.post('/refresh');

                        const newAccessToken = res.data?.access_token;
                        if (!newAccessToken) {
                            throw new Error("⚠ Refresh endpoint did not return access_token");
                        }

                        // Save new access token
                        setAuth(prev => ({
                            ...prev,
                            accessToken: newAccessToken
                        }));

                        // Apply token to queued requests
                        processQueue(null, newAccessToken);

                        // Retry original request with new token
                        originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
                        return axiosAuth(originalRequest);

                    } catch (refreshError) {

                        processQueue(refreshError, null);
                        setAuth(null);
                        localStorage.removeItem("user");
                        return Promise.reject(refreshError);

                    } finally {
                        isRefreshing = false;
                    }
                }

                return Promise.reject(error);
            }
        );

    }, [auth]);

    return axiosAuth;
};

export default useAxiosPrivate;
