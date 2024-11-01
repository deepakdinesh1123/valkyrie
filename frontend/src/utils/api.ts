import axios from 'axios';
import { Configuration, DefaultApi } from "@/api-client";

const config = new Configuration({
  basePath: `${import.meta.env.VITE_BASE_PATH}/api`,
});

const axiosInstance = axios.create({
  baseURL: config.basePath,
  timeout: 30000,

});

axiosInstance.interceptors.request.use(
  (config) => {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${localStorage.getItem('jwtToken') || import.meta.env.VITE_AUTH_TOKEN}`;
    return config;
  },
  (error) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error)
);

export const api = new DefaultApi(config, config.basePath, axiosInstance);