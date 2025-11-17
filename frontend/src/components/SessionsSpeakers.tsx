import { useEffect, useState } from 'react'
import { getSessions, getSpeakers } from '../api/endpoints'
import { SessionWithSpeakers } from '../types'

const SessionsSpeakers = () => {
  const [sessions, setSessions] = useState<SessionWithSpeakers[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [sessionsData, speakersData] = await Promise.all([
          getSessions(),
          getSpeakers(),
        ])

        // Map sessions with their speakers
        const sessionsWithSpeakers: SessionWithSpeakers[] = sessionsData.map((session) => ({
          ...session,
          speakers: speakersData.filter((speaker) =>
            session.speakerIds.includes(speaker.id || '')
          ),
        }))

        setSessions(sessionsWithSpeakers)
        setLoading(false)
      } catch (err) {
        setError('Failed to load sessions and speakers')
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) {
    return (
      <section id="sessions" className="py-16 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading sessions...</p>
          </div>
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section id="sessions" className="py-16 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center text-red-600">{error}</div>
        </div>
      </section>
    )
  }

  return (
    <section id="sessions" className="py-16 px-4 bg-white">
      <div className="max-w-6xl mx-auto">
        <h2 className="text-3xl md:text-4xl font-bold text-center mb-12 gradient-text">
          Sessions & Speakers
        </h2>
        {sessions.length === 0 ? (
          <div className="text-center text-gray-500 py-12">
            <p className="text-lg">No sessions available yet. Check back soon!</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {sessions.map((session) => (
              <div
                key={session.id}
                className="bg-white border border-gray-200 rounded-lg p-6 shadow-md hover:shadow-xl transition-all duration-300 transform hover:-translate-y-2"
              >
                <div className="mb-4">
                  <span className="text-sm font-semibold text-blue-600 bg-blue-50 px-3 py-1 rounded-full">
                    {session.time} • {session.duration}
                  </span>
                </div>
                <h3 className="text-xl font-bold mb-3 text-gray-900">{session.title}</h3>
                <p className="text-gray-600 mb-4 line-clamp-3">{session.description}</p>
                {session.speakers.length > 0 && (
                  <div className="border-t pt-4">
                    <p className="text-sm font-semibold text-gray-700 mb-2">Speakers:</p>
                    <div className="space-y-2">
                      {session.speakers.map((speaker) => (
                        <div key={speaker.id} className="flex items-center gap-2">
                          {speaker.imageUrl && (
                            <img
                              src={speaker.imageUrl}
                              alt={speaker.name}
                              className="w-10 h-10 rounded-full object-cover"
                            />
                          )}
                          <div>
                            <p className="font-medium text-gray-900">{speaker.name}</p>
                            {speaker.bio && (
                              <p className="text-xs text-gray-500 line-clamp-1">{speaker.bio}</p>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

export default SessionsSpeakers

