import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 60000, // Increase timeout to 60 seconds for large files
  headers: {
    'Content-Type': 'multipart/form-data',
  },
});

// Add response interceptor for better error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const customError = {
      message: 'Something went wrong. Please try again.',
      status: error.response?.status || 500,
      details: error.response?.data?.error || error.message,
    };

    if (error.code === 'ECONNABORTED') {
      customError.message = 'Upload timeout. Please try again.';
    }

    return Promise.reject(customError);
  }
);

export default api;
