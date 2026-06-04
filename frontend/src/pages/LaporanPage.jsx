import { useEffect, useState } from 'react'
import api from '../api/client'

export default function LaporanPage() {
  const [lap, setLap] = useState(null)

  useEffect(() => {
    api.get('/laporan').then((res) => setLap(res.data.data)).catch(() => {})
  }, [])

  if (!lap) return <p>Memuat laporan...</p>

  return (
    <div>
      <div className="page-head no-print-flex">
        <h2>Laporan Hasil Belajar</h2>
        <button className="btn btn-primary" onClick={() => window.print()}>🖨️ Cetak</button>
      </div>

      <div className="card-grid">
        <Stat label="Total Nilai" value={lap.total_nilai} />
        <Stat label="Lulus" value={lap.jumlah_lulus} cls="ok" />
        <Stat label="Tidak Lulus" value={lap.jumlah_tidak_lulus} cls="fail" />
        <Stat label="Rata-rata" value={lap.rata_rata_akhir} />
        <Stat label="Tertinggi" value={lap.nilai_tertinggi} />
        <Stat label="Terendah" value={lap.nilai_terendah} />
      </div>

      <div className="panel">
        <h3>Detail Nilai</h3>
        <table className="table">
          <thead>
            <tr><th>NIS</th><th>Nama</th><th>Kelas</th><th>Mapel</th><th>Nilai Akhir</th><th>Status</th></tr>
          </thead>
          <tbody>
            {(lap.detail || []).map((n) => (
              <tr key={n.id}>
                <td>{n.nis}</td>
                <td>{n.nama_siswa}</td>
                <td>{n.kelas}</td>
                <td>{n.mata_pelajaran}</td>
                <td><strong>{n.nilai_akhir}</strong></td>
                <td><span className={`tag ${n.status === 'LULUS' ? 'tag-pass' : 'tag-fail'}`}>{n.status}</span></td>
              </tr>
            ))}
            {(lap.detail || []).length === 0 && <tr><td colSpan="6" className="empty">Belum ada data</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function Stat({ label, value, cls }) {
  return (
    <div className={`stat-card ${cls ? 'stat-' + cls : ''}`}>
      <span className="stat-value">{value}</span>
      <span className="stat-label">{label}</span>
    </div>
  )
}
