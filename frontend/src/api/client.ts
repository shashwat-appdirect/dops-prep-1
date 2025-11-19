import axios from 'axios'

// Use relative URL for production (same domain), or env var for local development
const API_URL = import.meta.env.VITE_API_URL || ''

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

