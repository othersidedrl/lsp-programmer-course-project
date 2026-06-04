import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const demoAccounts = [
  { role: 'Admin', username: 'admin', password: 'admin123' },
  { role: 'Guru', username: 'guru', password: 'guru123' },
  { role: 'Siswa', username: 'siswa', password: 'siswa123' },
]

export default function Login() {
  const { login, user } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  // Jika sudah login, langsung arahkan ke beranda.
  useEffect(() => {
    if (user) navigate('/', { replace: true })
  }, [user, navigate])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(username, password)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err.response?.data?.message || 'Login gagal. Pastikan server backend berjalan.')
    } finally {
      setLoading(false)
    }
  }

  const fillDemo = (acc) => {
    setUsername(acc.username)
    setPassword(acc.password)
  }

  return (
    <div className="login-page">
      <form className="login-card" onSubmit={handleSubmit}>
        <div className="login-head">
          <div className="brand-icon">🎓</div>
          <h1>Pengolahan Nilai Siswa</h1>
          <p>Silakan masuk sesuai peran Anda</p>
        </div>

        {error && <div className="alert alert-error">{error}</div>}

        <label>Username</label>
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="masukkan username"
          autoFocus
        />

        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="masukkan password"
        />

        <button className="btn btn-primary btn-block" disabled={loading}>
          {loading ? 'Memproses...' : 'Masuk'}
        </button>

        <div className="demo-box">
          <small>Akun demo (klik untuk mengisi):</small>
          <div className="demo-list">
            {demoAccounts.map((a) => (
              <button
                type="button"
                key={a.username}
                className="demo-chip"
                onClick={() => fillDemo(a)}
              >
                <strong>{a.role}</strong> — {a.username} / {a.password}
              </button>
            ))}
          </div>
        </div>
      </form>
    </div>
  )
}
