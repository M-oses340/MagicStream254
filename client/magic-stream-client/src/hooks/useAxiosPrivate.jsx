import { useEffect } from "react";
import axios from "axios";
import useAuth from "./useAuth";

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {
  const { auth, setAuth } = useAuth();

  // Axios instance
  const axiosPrivate = axios.create({
    baseURL: apiUrl,
    withCredentials: true, // required for cookies
  });

  let isRefreshing = false;
  let failedQueue = [];

  const processQueue = (error, token = null) => {
    failedQueue.forEach((prom) => {
      if (error) {
        prom.reject(error);
      } else {
        prom.resolve(token);
      }
    });
    failedQueue = [];
  };

  useEffect(() => {
    // =========================
    // REQUEST INTERCEPTOR
    // =========================
    const requestInterceptor = axiosPrivate.interceptors.request.use(
      (config) => {
        if (auth?.accessToken) {
          config.headers["Authorization"] = `Bearer ${auth.accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // =========================
    // RESPONSE INTERCEPTOR
    // =========================
    const responseInterceptor = axiosPrivate.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        // If refresh endpoint itself fails → logout
        if (originalRequest.url.includes("/refresh") && error?.response?.status === 401) {
          setAuth(null);
          localStorage.removeItem("user");
          return Promise.reject(error);
        }

        // Handle 401 errors (access token expired)
        if (error?.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              failedQueue.push({ resolve, reject });
            })
              .then((token) => {
                originalRequest.headers["Authorization"] = `Bearer ${token}`;
                return axiosPrivate(originalRequest);
              })
              .catch((err) => Promise.reject(err));
          }

          isRefreshing = true;

          try {
            const res = await axiosPrivate.post(
              "/refresh",
              {},
              { withCredentials: true } // ✅ important to send refresh cookie
            );

            const newAccessToken = res.data?.access_token;
            if (!newAccessToken)
              throw new Error("⚠ Refresh endpoint did not return access_token");

            // Update auth state
            setAuth((prev) => ({ ...prev, accessToken: newAccessToken }));

            // Retry queued requests
            processQueue(null, newAccessToken);

            // Retry original request
            originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
            return axiosPrivate(originalRequest);
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

    // Cleanup interceptors on unmount
    return () => {
      axiosPrivate.interceptors.request.eject(requestInterceptor);
      axiosPrivate.interceptors.response.eject(responseInterceptor);
    };
  }, [auth]);

  return axiosPrivate;
};

export default useAxiosPrivate;
