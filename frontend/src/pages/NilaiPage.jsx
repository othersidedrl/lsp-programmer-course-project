import { useEffect, useState } from 'react'
import api from '../api/client'
import { useAuth } from '../context/AuthContext'

const emptyForm = { nis: '', mata_pelajaran: '', id_guru: '', tugas: '', uts: '', uas: '' }

export default function NilaiPage() {
  const { user } = useAuth()
  const isAdmin = user.role === 'admin'

  const [list, setList] = useState([])
  const [siswaList, setSiswaList] = useState([])
  const [guruList, setGuruList] = useState([])
  const [form, setForm] = useState(emptyForm)
  const [editingId, setEditingId] = useState(null)
  const [msg, setMsg] = useState(null)

  const notify = (type, text) => {
    setMsg({ type, text })
    setTimeout(() => setMsg(null), 3500)
  }

  const load = async () => {
    const [n, s, g] = await Promise.all([
      api.get('/nilai'),
      api.get('/siswa'),
      api.get('/guru'),
    ])
    setList(n.data.data || [])
    setSiswaList(s.data.data || [])
    setGuruList(g.data.data || [])
  }
  useEffect(() => { load() }, [])

  // Mata pelajaran yang diampu guru yang sedang login.
  const guruMapel = !isAdmin
    ? guruList.find((g) => g.id_guru === user.ref_id)?.mata_pelajaran || ''
    : ''

  const submit = async (e) => {
    e.preventDefault()
    const payload = {
      nis: form.nis,
      mata_pelajaran: isAdmin ? form.mata_pelajaran : guruMapel,
      id_guru: isAdmin ? form.id_guru : user.ref_id,
      tugas: Number(form.tugas),
      uts: Number(form.uts),
      uas: Number(form.uas),
    }
    try {
      if (editingId) {
        await api.put(`/nilai/${editingId}`, payload)
        notify('success', 'Nilai diperbarui — status validasi direset, perlu validasi ulang')
      } else {
        await api.post('/nilai', payload)
        notify('success', 'Nilai disimpan & nilai akhir dihitung otomatis')
      }
      setForm(emptyForm)
      setEditingId(null)
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menyimpan nilai')
    }
  }

  const startEdit = (n) => {
    setEditingId(n.id)
    setForm({
      nis: n.nis,
      mata_pelajaran: n.mata_pelajaran,
      id_guru: n.id_guru,
      tugas: n.tugas,
      uts: n.uts,
      uas: n.uas,
    })
  }
  const cancel = () => { setForm(emptyForm); setEditingId(null) }

  const validate = async (id) => {
    try {
      await api.patch(`/nilai/${id}/validasi`)
      notify('success', 'Nilai berhasil divalidasi')
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal memvalidasi')
    }
  }

  const remove = async (id) => {
    if (!window.confirm('Hapus data nilai ini?')) return
    try {
      await api.delete(`/nilai/${id}`)
      notify('success', 'Nilai berhasil dihapus')
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menghapus')
    }
  }

  return (
    <div>
      <div className="page-head">
        <h2>{isAdmin ? 'Data Nilai' : 'Input & Rekap Nilai'}</h2>
        {!isAdmin && <p>Mata pelajaran Anda: <strong>{guruMapel || '—'}</strong></p>}
      </div>
      {msg && <div className={`alert alert-${msg.type}`}>{msg.text}</div>}

      <div className="panel">
        <h3>{editingId ? 'Ubah Nilai' : 'Input Nilai Baru'}</h3>
        <form onSubmit={submit} className="form-row">
          <div>
            <label>Siswa</label>
            <select
              value={form.nis}
              disabled={!!editingId}
              onChange={(e) => setForm({ ...form, nis: e.target.value })}
              required
            >
              <option value="">-- pilih siswa --</option>
              {siswaList.map((s) => (
                <option key={s.nis} value={s.nis}>{s.nis} - {s.nama}</option>
              ))}
            </select>
          </div>

          {isAdmin ? (
            <>
              <div>
                <label>Mata Pelajaran</label>
                <input
                  value={form.mata_pelajaran}
                  onChange={(e) => setForm({ ...form, mata_pelajaran: e.target.value })}
                  required
                />
              </div>
              <div>
                <label>Guru</label>
                <select
                  value={form.id_guru}
                  onChange={(e) => setForm({ ...form, id_guru: e.target.value })}
                >
                  <option value="">-- pilih guru --</option>
                  {guruList.map((g) => (
                    <option key={g.id_guru} value={g.id_guru}>{g.nama} ({g.mata_pelajaran})</option>
                  ))}
                </select>
              </div>
            </>
          ) : (
            <div>
              <label>Mata Pelajaran</label>
              <input value={guruMapel} disabled />
            </div>
          )}

          <div>
            <label>Tugas</label>
            <input type="number" min="0" max="100" step="0.01" value={form.tugas}
              onChange={(e) => setForm({ ...form, tugas: e.target.value })} required />
          </div>
          <div>
            <label>UTS</label>
            <input type="number" min="0" max="100" step="0.01" value={form.uts}
              onChange={(e) => setForm({ ...form, uts: e.target.value })} required />
          </div>
          <div>
            <label>UAS</label>
            <input type="number" min="0" max="100" step="0.01" value={form.uas}
              onChange={(e) => setForm({ ...form, uas: e.target.value })} required />
          </div>

          <div className="form-actions">
            <button className="btn btn-primary">{editingId ? 'Simpan' : 'Hitung & Simpan'}</button>
            {editingId && <button type="button" className="btn btn-ghost" onClick={cancel}>Batal</button>}
          </div>
        </form>
        <small className="hint">
          Nilai akhir dihitung otomatis: 30% Tugas + 30% UTS + 40% UAS. Rentang nilai valid 0–100.
        </small>
      </div>

      <div className="panel">
        <h3>Rekap Nilai ({list.length})</h3>
        <table className="table">
          <thead>
            <tr>
              <th>Nama</th><th>Kelas</th><th>Mapel</th>
              <th>Tugas</th><th>UTS</th><th>UAS</th>
              <th>Akhir</th><th>Status</th><th>Validasi</th><th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {list.map((n) => (
              <tr key={n.id}>
                <td>{n.nama_siswa}</td>
                <td>{n.kelas}</td>
                <td>{n.mata_pelajaran}</td>
                <td>{n.tugas}</td>
                <td>{n.uts}</td>
                <td>{n.uas}</td>
                <td><strong>{n.nilai_akhir}</strong></td>
                <td><span className={`tag ${n.status === 'LULUS' ? 'tag-pass' : 'tag-fail'}`}>{n.status}</span></td>
                <td>
                  {n.tervalidasi
                    ? <span className="tag tag-ok">✔ Tervalidasi</span>
                    : <span className="tag tag-warn">Belum</span>}
                </td>
                <td className="row-actions">
                  <button className="btn btn-sm" onClick={() => startEdit(n)}>Ubah</button>
                  {!n.tervalidasi && <button className="btn btn-sm btn-ok" onClick={() => validate(n.id)}>Validasi</button>}
                  {isAdmin && <button className="btn btn-sm btn-danger" onClick={() => remove(n.id)}>Hapus</button>}
                </td>
              </tr>
            ))}
            {list.length === 0 && <tr><td colSpan="10" className="empty">Belum ada data nilai</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
