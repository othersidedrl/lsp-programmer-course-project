import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

// Menu navigasi berbeda untuk tiap peran (hak akses pengguna).
const menuByRole = {
  admin: [
    { to: '/', label: 'Dashboard', end: true },
    { to: '/siswa', label: 'Data Siswa' },
    { to: '/guru', label: 'Data Guru' },
    { to: '/nilai', label: 'Data Nilai' },
    { to: '/laporan', label: 'Laporan' },
  ],
  guru: [
    { to: '/', label: 'Dashboard', end: true },
    { to: '/nilai', label: 'Input & Rekap Nilai' },
    { to: '/laporan', label: 'Laporan' },
  ],
  siswa: [
    { to: '/', label: 'Dashboard', end: true },
    { to: '/nilai-saya', label: 'Nilai Saya' },
  ],
}

export default function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const menu = menuByRole[user.role] || []

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-icon">🎓</span>
          <div>
            <strong>Nilai Siswa</strong>
            <small>Sistem Pengolahan Nilai</small>
          </div>
        </div>
        <nav>
          {menu.map((m) => (
            <NavLink key={m.to} to={m.to} end={m.end} className="nav-link">
              {m.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className="main">
        <header className="topbar">
          <span className={`role-badge role-${user.role}`}>{user.role.toUpperCase()}</span>
          <div className="topbar-right">
            <span>👤 {user.username}</span>
            <button className="btn btn-ghost" onClick={handleLogout}>Keluar</button>
          </div>
        </header>
        <main className="content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
