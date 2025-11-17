import apiClient from './client'
import { Registration, Speaker, Session, DesignationBreakdown } from '../types'

export const register = async (data: Omit<Registration, 'id' | 'createdAt'>) => {
  const response = await apiClient.post('/api/register', data)
  return response.data
}

export const getRegistrationCount = async (): Promise<number> => {
  const response = await apiClient.get('/api/registrations/count')
  return response.data.count
}

export const adminLogin = async (password: string): Promise<string> => {
  const response = await apiClient.post('/api/admin/login', { password })
  return response.data.token
}

export const getAttendees = async (): Promise<Registration[]> => {
  const response = await apiClient.get('/api/admin/attendees')
  return response.data
}

export const getAttendee = async (id: string): Promise<Registration> => {
  const response = await apiClient.get(`/api/admin/attendees/${id}`)
  return response.data
}

export const getSpeakers = async (): Promise<Speaker[]> => {
  // Use public endpoint for home screen, admin endpoint for admin dashboard
  const token = localStorage.getItem('admin_token')
  const endpoint = token ? '/api/admin/speakers' : '/api/speakers'
  const response = await apiClient.get(endpoint)
  return response.data || []
}

export const createSpeaker = async (data: Omit<Speaker, 'id'>): Promise<Speaker> => {
  const response = await apiClient.post('/api/admin/speakers', data)
  return response.data
}

export const updateSpeaker = async (id: string, data: Partial<Speaker>): Promise<Speaker> => {
  const response = await apiClient.put(`/api/admin/speakers/${id}`, data)
  return response.data
}

export const deleteSpeaker = async (id: string): Promise<void> => {
  await apiClient.delete(`/api/admin/speakers/${id}`)
}

export const getSessions = async (): Promise<Session[]> => {
  // Use public endpoint for home screen, admin endpoint for admin dashboard
  const token = localStorage.getItem('admin_token')
  const endpoint = token ? '/api/admin/sessions' : '/api/sessions'
  const response = await apiClient.get(endpoint)
  return response.data || []
}

export const createSession = async (data: Omit<Session, 'id'>): Promise<Session> => {
  const response = await apiClient.post('/api/admin/sessions', data)
  return response.data
}

export const updateSession = async (id: string, data: Partial<Session>): Promise<Session> => {
  const response = await apiClient.put(`/api/admin/sessions/${id}`, data)
  return response.data
}

export const deleteSession = async (id: string): Promise<void> => {
  await apiClient.delete(`/api/admin/sessions/${id}`)
}

export const getDesignationBreakdown = async (): Promise<DesignationBreakdown[]> => {
  const response = await apiClient.get('/api/admin/analytics/designations')
  return response.data
}

