import axios from 'axios'

// Klien HTTP terpusat. baseURL "/api" akan di-proxy ke backend Go saat dev.
const api = axios.create({ baseURL: '/api' })

// Sisipkan token JWT pada setiap request bila tersedia.
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Bila token kedaluwarsa/invalid (401), bersihkan sesi & arahkan ke login.
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  },
)

export default api
