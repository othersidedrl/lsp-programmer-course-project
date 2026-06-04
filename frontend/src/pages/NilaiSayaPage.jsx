import { useEffect, useState } from 'react'
import api from '../api/client'
import { useAuth } from '../context/AuthContext'

export default function NilaiSayaPage() {
  const { user } = useAuth()
  const [list, setList] = useState([])

  useEffect(() => {
    api.get('/nilai').then((res) => setList(res.data.data || [])).catch(() => {})
  }, [])

  const lulus = list.filter((n) => n.status === 'LULUS').length
  const semuaLulus = list.length > 0 && lulus === list.length

  return (
    <div>
      <div className="page-head">
        <h2>Nilai Saya</h2>
        <p>NIS: <strong>{user.ref_id}</strong></p>
      </div>

      {list.length > 0 && (
        <div className={`banner ${semuaLulus ? 'banner-ok' : 'banner-warn'}`}>
          {semuaLulus
            ? '🎉 Selamat! Anda dinyatakan LULUS pada semua mata pelajaran.'
            : `Anda lulus pada ${lulus} dari ${list.length} mata pelajaran.`}
        </div>
      )}

      <div className="panel">
        <h3>Daftar Nilai</h3>
        <table className="table">
          <thead>
            <tr><th>Mata Pelajaran</th><th>Tugas</th><th>UTS</th><th>UAS</th><th>Nilai Akhir</th><th>Status</th></tr>
          </thead>
          <tbody>
            {list.map((n) => (
              <tr key={n.id}>
                <td>{n.mata_pelajaran}</td>
                <td>{n.tugas}</td>
                <td>{n.uts}</td>
                <td>{n.uas}</td>
                <td><strong>{n.nilai_akhir}</strong></td>
                <td><span className={`tag ${n.status === 'LULUS' ? 'tag-pass' : 'tag-fail'}`}>{n.status}</span></td>
              </tr>
            ))}
            {list.length === 0 && <tr><td colSpan="6" className="empty">Belum ada nilai yang tersedia</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
