import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { adminLogin } from '../api/endpoints'

const Footer = () => {
  const [showLogin, setShowLogin] = useState(false)
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const token = await adminLogin(password)
      localStorage.setItem('admin_token', token)
      setShowLogin(false)
      navigate('/admin')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <footer className="bg-gray-900 text-white py-8 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <p className="text-gray-400">
              &copy; {new Date().getFullYear()} AppDirect India AI Workshop. All rights reserved.
            </p>
          </div>
          <button
            onClick={() => setShowLogin(true)}
            className="text-sm text-gray-400 hover:text-white transition-colors"
          >
            Admin Login
          </button>
        </div>
      </div>

      {/* Login Modal */}
      {showLogin && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 animate-fade-in">
          <div className="bg-white rounded-lg p-8 max-w-md w-full animate-slide-down">
            <h3 className="text-2xl font-bold text-gray-900 mb-4">Admin Login</h3>
            <form onSubmit={handleLogin} className="space-y-4">
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                  Password
                </label>
                <input
                  type="password"
                  id="password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Enter admin password"
                />
              </div>
              {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-2 rounded-lg text-sm">
                  {error}
                </div>
              )}
              <div className="flex gap-3">
                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 gradient-bg text-white py-2 rounded-lg font-semibold hover:opacity-90 transition-all disabled:opacity-50"
                >
                  {loading ? 'Logging in...' : 'Login'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowLogin(false)
                    setPassword('')
                    setError(null)
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-all"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </footer>
  )
}

export default Footer

