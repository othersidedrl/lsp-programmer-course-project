import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../api/client'
import { useAuth } from '../context/AuthContext'

export default function Dashboard() {
  const { user } = useAuth()
  const [stats, setStats] = useState({})

  useEffect(() => {
    async function load() {
      try {
        if (user.role === 'admin') {
          const [s, g, n] = await Promise.all([
            api.get('/siswa'),
            api.get('/guru'),
            api.get('/nilai'),
          ])
          setStats({
            siswa: s.data.data.length,
            guru: g.data.data.length,
            nilai: n.data.data.length,
          })
        } else {
          const n = await api.get('/nilai')
          const data = n.data.data || []
          setStats({
            nilai: data.length,
            lulus: data.filter((x) => x.status === 'LULUS').length,
          })
        }
      } catch {
        // biarkan kartu menampilkan tanda "—"
      }
    }
    load()
  }, [user.role])

  return (
    <div>
      <div className="page-head">
        <h2>Dashboard</h2>
        <p>Selamat datang, <strong>{user.username}</strong> — peran <strong>{user.role}</strong>.</p>
      </div>

      <div className="card-grid">
        {user.role === 'admin' && (
          <>
            <StatCard label="Total Siswa" value={stats.siswa ?? '—'} to="/siswa" />
            <StatCard label="Total Guru" value={stats.guru ?? '—'} to="/guru" />
            <StatCard label="Total Nilai" value={stats.nilai ?? '—'} to="/nilai" />
            <StatCard label="Lihat Laporan" value="📊" to="/laporan" />
          </>
        )}
        {user.role === 'guru' && (
          <>
            <StatCard label="Nilai Diampu" value={stats.nilai ?? '—'} to="/nilai" />
            <StatCard label="Jumlah Lulus" value={stats.lulus ?? '—'} to="/laporan" />
            <StatCard label="Input Nilai" value="➕" to="/nilai" />
          </>
        )}
        {user.role === 'siswa' && (
          <>
            <StatCard label="Mata Pelajaran" value={stats.nilai ?? '—'} to="/nilai-saya" />
            <StatCard label="Mapel Lulus" value={stats.lulus ?? '—'} to="/nilai-saya" />
          </>
        )}
      </div>

      <div className="info-panel">
        <h3>Tentang Aplikasi</h3>
        <p>
          Aplikasi ini membantu pengelolaan nilai siswa: input nilai, perhitungan
          nilai akhir otomatis <code>(30% Tugas + 30% UTS + 40% UAS)</code>,
          penentuan status kelulusan (<strong>nilai akhir ≥ 70 dinyatakan LULUS</strong>),
          serta penyajian laporan hasil belajar.
        </p>
      </div>
    </div>
  )
}

function StatCard({ label, value, to }) {
  return (
    <Link to={to} className="stat-card">
      <span className="stat-value">{value}</span>
      <span className="stat-label">{label}</span>
    </Link>
  )
}
