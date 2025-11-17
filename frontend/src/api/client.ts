import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests if available (only for admin endpoints)
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  // Only add token for admin endpoints
  if (token && config.url && config.url.startsWith('/api/admin')) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default apiClient

