import { useState, useEffect } from 'react'
import { register, getRegistrationCount } from '../api/endpoints'

const DESIGNATIONS = [
  'Software Engineer',
  'Product Manager',
  'Designer',
  'Data Scientist',
  'Student',
  'Other',
]

const RegistrationForm = () => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    designation: DESIGNATIONS[0],
  })
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(false)
  const [showSuccess, setShowSuccess] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchCount = async () => {
      try {
        const registrationCount = await getRegistrationCount()
        setCount(registrationCount)
      } catch (err) {
        console.error('Failed to fetch registration count:', err)
      }
    }

    fetchCount()
    // Refresh count every 30 seconds
    const interval = setInterval(fetchCount, 30000)
    return () => clearInterval(interval)
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      await register(formData)
      setShowSuccess(true)
      setFormData({ name: '', email: '', designation: DESIGNATIONS[0] })
      // Refresh count after successful registration
      const registrationCount = await getRegistrationCount()
      setCount(registrationCount)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to register. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section id="register" className="py-16 px-4 bg-gray-50">
      <div className="max-w-2xl mx-auto">
        <div className="bg-white rounded-lg shadow-lg p-8">
          <div className="text-center mb-8">
            <h2 className="text-3xl md:text-4xl font-bold mb-4 gradient-text">
              Register for the Workshop
            </h2>
            <div className="inline-flex items-center gap-2 bg-blue-50 text-blue-700 px-4 py-2 rounded-full">
              <span className="text-2xl font-bold">{count}</span>
              <span className="text-sm">people registered</span>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                Full Name *
              </label>
              <input
                type="text"
                id="name"
                required
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Enter your full name"
              />
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                Email Address *
              </label>
              <input
                type="email"
                id="email"
                required
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="your.email@example.com"
              />
            </div>

            <div>
              <label
                htmlFor="designation"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Designation *
              </label>
              <select
                id="designation"
                required
                value={formData.designation}
                onChange={(e) => setFormData({ ...formData, designation: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all bg-white"
              >
                {DESIGNATIONS.map((designation) => (
                  <option key={designation} value={designation}>
                    {designation}
                  </option>
                ))}
              </select>
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full gradient-bg text-white py-3 rounded-lg font-semibold hover:opacity-90 transition-all duration-300 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Registering...' : 'Register Now'}
            </button>
          </form>
        </div>
      </div>

      {/* Success Modal */}
      {showSuccess && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 animate-fade-in">
          <div className="bg-white rounded-lg p-8 max-w-md w-full animate-slide-down">
            <div className="text-center">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg
                  className="w-8 h-8 text-green-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">Registration Successful!</h3>
              <p className="text-gray-600 mb-6">
                Thank you for registering. We'll send you a confirmation email shortly.
              </p>
              <button
                onClick={() => setShowSuccess(false)}
                className="gradient-bg text-white px-6 py-2 rounded-lg font-semibold hover:opacity-90 transition-all"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </section>
  )
}

export default RegistrationForm

