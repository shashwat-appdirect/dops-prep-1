import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  getAttendees,
  getSpeakers,
  getSessions,
  getDesignationBreakdown,
  createSpeaker,
  updateSpeaker,
  deleteSpeaker,
  createSession,
  updateSession,
  deleteSession,
} from '../api/endpoints'
import { Registration, Speaker, Session, DesignationBreakdown } from '../types'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'

const COLORS = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#ef4444']

const AdminDashboard = () => {
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<'attendees' | 'speakers' | 'sessions' | 'analytics'>(
    'attendees'
  )
  const [attendees, setAttendees] = useState<Registration[]>([])
  const [speakers, setSpeakers] = useState<Speaker[]>([])
  const [sessions, setSessions] = useState<Session[]>([])
  const [breakdown, setBreakdown] = useState<DesignationBreakdown[]>([])
  const [loading, setLoading] = useState(true)
  const [showSpeakerModal, setShowSpeakerModal] = useState(false)
  const [showSessionModal, setShowSessionModal] = useState(false)
  const [editingSpeaker, setEditingSpeaker] = useState<Speaker | null>(null)
  const [editingSession, setEditingSession] = useState<Session | null>(null)
  const [speakerForm, setSpeakerForm] = useState<Omit<Speaker, 'id'>>({
    name: '',
    bio: '',
    imageUrl: '',
    linkedinUrl: '',
    twitterUrl: '',
  })
  const [sessionForm, setSessionForm] = useState<Omit<Session, 'id'>>({
    title: '',
    description: '',
    time: '',
    duration: '',
    speakerIds: [],
  })

  useEffect(() => {
    const token = localStorage.getItem('admin_token')
    if (!token) {
      navigate('/')
      return
    }
    loadData()
  }, [navigate])

  const loadData = async () => {
    try {
      setLoading(true)
      const [attendeesData, speakersData, sessionsData, breakdownData] = await Promise.all([
        getAttendees(),
        getSpeakers(),
        getSessions(),
        getDesignationBreakdown(),
      ])
      setAttendees(Array.isArray(attendeesData) ? attendeesData : [])
      setSpeakers(Array.isArray(speakersData) ? speakersData : [])
      setSessions(Array.isArray(sessionsData) ? sessionsData : [])
      setBreakdown(Array.isArray(breakdownData) ? breakdownData : [])
    } catch (err) {
      console.error('Failed to load data:', err)
      // Don't logout on error, just set empty arrays
      setAttendees([])
      setSpeakers([])
      setSessions([])
      setBreakdown([])
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('admin_token')
    navigate('/')
  }

  const handleSpeakerSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingSpeaker) {
        await updateSpeaker(editingSpeaker.id!, speakerForm)
      } else {
        await createSpeaker(speakerForm)
      }
      await loadData()
      setShowSpeakerModal(false)
      resetSpeakerForm()
    } catch (err) {
      console.error('Failed to save speaker:', err)
    }
  }

  const handleSessionSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingSession) {
        await updateSession(editingSession.id!, sessionForm)
      } else {
        await createSession(sessionForm)
      }
      await loadData()
      setShowSessionModal(false)
      resetSessionForm()
    } catch (err) {
      console.error('Failed to save session:', err)
    }
  }

  const handleDeleteSpeaker = async (id: string) => {
    if (confirm('Are you sure you want to delete this speaker?')) {
      try {
        await deleteSpeaker(id)
        await loadData()
      } catch (err) {
        console.error('Failed to delete speaker:', err)
      }
    }
  }

  const handleDeleteSession = async (id: string) => {
    if (confirm('Are you sure you want to delete this session?')) {
      try {
        await deleteSession(id)
        await loadData()
      } catch (err) {
        console.error('Failed to delete session:', err)
      }
    }
  }

  const resetSpeakerForm = () => {
    setSpeakerForm({ name: '', bio: '', imageUrl: '', linkedinUrl: '', twitterUrl: '' })
    setEditingSpeaker(null)
  }

  const resetSessionForm = () => {
    setSessionForm({ title: '', description: '', time: '', duration: '', speakerIds: [] })
    setEditingSession(null)
  }

  const openSpeakerModal = (speaker?: Speaker) => {
    if (speaker) {
      setEditingSpeaker(speaker)
      setSpeakerForm({
        name: speaker.name,
        bio: speaker.bio,
        imageUrl: speaker.imageUrl || '',
        linkedinUrl: speaker.linkedinUrl || '',
        twitterUrl: speaker.twitterUrl || '',
      })
    } else {
      resetSpeakerForm()
    }
    setShowSpeakerModal(true)
  }

  const openSessionModal = (session?: Session) => {
    if (session) {
      setEditingSession(session)
      setSessionForm({
        title: session.title,
        description: session.description,
        time: session.time,
        duration: session.duration,
        speakerIds: session.speakerIds,
      })
    } else {
      resetSessionForm()
    }
    setShowSessionModal(true)
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold gradient-text">Admin Dashboard</h1>
          <button
            onClick={handleLogout}
            className="text-red-600 hover:text-red-700 font-medium"
          >
            Logout
          </button>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-6">
        {/* Tabs */}
        <div className="bg-white rounded-lg shadow-sm mb-6">
          <div className="flex border-b">
            {(['attendees', 'speakers', 'sessions', 'analytics'] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-6 py-3 font-medium capitalize transition-colors ${
                  activeTab === tab
                    ? 'text-blue-600 border-b-2 border-blue-600'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {tab}
              </button>
            ))}
          </div>
        </div>

        {/* Attendees Tab */}
        {activeTab === 'attendees' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-bold mb-4">Attendees ({attendees.length})</h2>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-semibold">Name</th>
                    <th className="text-left py-3 px-4 font-semibold">Email</th>
                    <th className="text-left py-3 px-4 font-semibold">Designation</th>
                    <th className="text-left py-3 px-4 font-semibold">Registered</th>
                  </tr>
                </thead>
                <tbody>
                  {attendees.map((attendee) => (
                    <tr key={attendee.id} className="border-b hover:bg-gray-50">
                      <td className="py-3 px-4">{attendee.name}</td>
                      <td className="py-3 px-4">{attendee.email}</td>
                      <td className="py-3 px-4">{attendee.designation}</td>
                      <td className="py-3 px-4">
                        {attendee.createdAt
                          ? new Date(attendee.createdAt).toLocaleDateString()
                          : '-'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Speakers Tab */}
        {activeTab === 'speakers' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">Speakers ({speakers.length})</h2>
              <button
                onClick={() => openSpeakerModal()}
                className="gradient-bg text-white px-4 py-2 rounded-lg font-semibold hover:opacity-90"
              >
                Add Speaker
              </button>
            </div>
            {speakers.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500 text-lg mb-4">No speakers added yet.</p>
                <p className="text-gray-400 text-sm">Click "Add Speaker" to get started.</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {speakers.map((speaker) => (
                  <div key={speaker.id} className="border rounded-lg p-4">
                    {speaker.imageUrl && (
                      <img
                        src={speaker.imageUrl}
                        alt={speaker.name}
                        className="w-full h-48 object-cover rounded-lg mb-3"
                      />
                    )}
                    <h3 className="font-bold text-lg mb-2">{speaker.name}</h3>
                    <p className="text-gray-600 text-sm mb-3 line-clamp-3">{speaker.bio}</p>
                    <div className="flex gap-2">
                      <button
                        onClick={() => openSpeakerModal(speaker)}
                        className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDeleteSpeaker(speaker.id!)}
                        className="text-red-600 hover:text-red-700 text-sm font-medium"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Sessions Tab */}
        {activeTab === 'sessions' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">Sessions ({sessions.length})</h2>
              <button
                onClick={() => openSessionModal()}
                className="gradient-bg text-white px-4 py-2 rounded-lg font-semibold hover:opacity-90"
              >
                Add Session
              </button>
            </div>
            {sessions.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500 text-lg mb-4">No sessions added yet.</p>
                <p className="text-gray-400 text-sm">Click "Add Session" to get started.</p>
              </div>
            ) : (
              <div className="space-y-4">
                {sessions.map((session) => (
                  <div key={session.id} className="border rounded-lg p-4">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <span className="text-sm font-semibold text-blue-600 bg-blue-50 px-2 py-1 rounded">
                            {session.time} • {session.duration}
                          </span>
                        </div>
                        <h3 className="font-bold text-lg mb-2">{session.title}</h3>
                        <p className="text-gray-600 mb-2">{session.description}</p>
                        <p className="text-sm text-gray-500">
                          Speakers: {session.speakerIds.length}
                        </p>
                      </div>
                      <div className="flex gap-2 ml-4">
                        <button
                          onClick={() => openSessionModal(session)}
                          className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteSession(session.id!)}
                          className="text-red-600 hover:text-red-700 text-sm font-medium"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Analytics Tab */}
        {activeTab === 'analytics' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-bold mb-6">Designation Breakdown</h2>
            {breakdown.length > 0 ? (
              <div className="max-w-2xl mx-auto">
                <ResponsiveContainer width="100%" height={400}>
                  <PieChart>
                    <Pie
                      data={breakdown}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ designation, count }) => `${designation}: ${count}`}
                      outerRadius={120}
                      fill="#8884d8"
                      dataKey="count"
                    >
                      {breakdown.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
                <div className="mt-6 space-y-2">
                  {breakdown.map((item, index) => (
                    <div key={item.designation} className="flex items-center gap-3">
                      <div
                        className="w-4 h-4 rounded"
                        style={{ backgroundColor: COLORS[index % COLORS.length] }}
                      ></div>
                      <span className="font-medium">{item.designation}</span>
                      <span className="text-gray-600">({item.count})</span>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-8">No data available</p>
            )}
          </div>
        )}
      </div>

      {/* Speaker Modal */}
      {showSpeakerModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">
              {editingSpeaker ? 'Edit Speaker' : 'Add Speaker'}
            </h3>
            <form onSubmit={handleSpeakerSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Name *</label>
                <input
                  type="text"
                  required
                  value={speakerForm.name}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, name: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Bio *</label>
                <textarea
                  required
                  value={speakerForm.bio}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, bio: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                  rows={4}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Image URL</label>
                <input
                  type="url"
                  value={speakerForm.imageUrl}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, imageUrl: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">LinkedIn URL</label>
                <input
                  type="url"
                  value={speakerForm.linkedinUrl}
                  onChange={(e) =>
                    setSpeakerForm({ ...speakerForm, linkedinUrl: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Twitter URL</label>
                <input
                  type="url"
                  value={speakerForm.twitterUrl}
                  onChange={(e) =>
                    setSpeakerForm({ ...speakerForm, twitterUrl: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded-lg"
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="gradient-bg text-white px-4 py-2 rounded-lg font-semibold"
                >
                  Save
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowSpeakerModal(false)
                    resetSpeakerForm()
                  }}
                  className="px-4 py-2 border rounded-lg"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Session Modal */}
      {showSessionModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">
              {editingSession ? 'Edit Session' : 'Add Session'}
            </h3>
            <form onSubmit={handleSessionSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Title *</label>
                <input
                  type="text"
                  required
                  value={sessionForm.title}
                  onChange={(e) => setSessionForm({ ...sessionForm, title: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Description *</label>
                <textarea
                  required
                  value={sessionForm.description}
                  onChange={(e) =>
                    setSessionForm({ ...sessionForm, description: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded-lg"
                  rows={4}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Time *</label>
                  <input
                    type="text"
                    required
                    value={sessionForm.time}
                    onChange={(e) => setSessionForm({ ...sessionForm, time: e.target.value })}
                    className="w-full px-3 py-2 border rounded-lg"
                    placeholder="e.g., 10:00 AM"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Duration *</label>
                  <input
                    type="text"
                    required
                    value={sessionForm.duration}
                    onChange={(e) => setSessionForm({ ...sessionForm, duration: e.target.value })}
                    className="w-full px-3 py-2 border rounded-lg"
                    placeholder="e.g., 45 min"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Speakers</label>
                <select
                  multiple
                  value={sessionForm.speakerIds}
                  onChange={(e) => {
                    const selectedIds = Array.from(e.target.selectedOptions, (option) => option.value)
                    setSessionForm({
                      ...sessionForm,
                      speakerIds: selectedIds,
                    })
                  }}
                  className="w-full px-3 py-2 border rounded-lg min-h-[120px]"
                  size={Math.min(speakers.length || 1, 5)}
                >
                  {speakers.length === 0 ? (
                    <option value="" disabled>
                      No speakers available. Add speakers first.
                    </option>
                  ) : (
                    speakers.map((speaker) => (
                      <option key={speaker.id} value={speaker.id || ''}>
                        {speaker.name}
                      </option>
                    ))
                  )}
                </select>
                <p className="text-xs text-gray-500 mt-1">
                  {sessionForm.speakerIds.length > 0
                    ? `${sessionForm.speakerIds.length} speaker(s) selected. Hold Ctrl/Cmd to select multiple.`
                    : 'Hold Ctrl/Cmd (Mac) or Ctrl (Windows) to select multiple speakers.'}
                </p>
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="gradient-bg text-white px-4 py-2 rounded-lg font-semibold"
                >
                  Save
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowSessionModal(false)
                    resetSessionForm()
                  }}
                  className="px-4 py-2 border rounded-lg"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default AdminDashboard

