import { useEffect, useState } from 'react'
import api from '../api/client'

const empty = { nis: '', nama: '', kelas: '' }

export default function SiswaPage() {
  const [list, setList] = useState([])
  const [form, setForm] = useState(empty)
  const [editing, setEditing] = useState(false)
  const [msg, setMsg] = useState(null)
  const [loading, setLoading] = useState(false)

  const notify = (type, text) => {
    setMsg({ type, text })
    setTimeout(() => setMsg(null), 3000)
  }

  const load = async () => {
    setLoading(true)
    try {
      const res = await api.get('/siswa')
      setList(res.data.data || [])
    } finally {
      setLoading(false)
    }
  }
  useEffect(() => { load() }, [])

  const submit = async (e) => {
    e.preventDefault()
    try {
      if (editing) {
        await api.put(`/siswa/${form.nis}`, form)
        notify('success', 'Data siswa berhasil diperbarui')
      } else {
        await api.post('/siswa', form)
        notify('success', 'Siswa berhasil ditambahkan')
      }
      setForm(empty)
      setEditing(false)
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menyimpan data')
    }
  }

  const edit = (s) => { setForm(s); setEditing(true) }
  const cancel = () => { setForm(empty); setEditing(false) }

  const remove = async (nis) => {
    if (!window.confirm(`Hapus siswa ${nis}? Nilai terkait juga ikut terhapus.`)) return
    try {
      await api.delete(`/siswa/${nis}`)
      notify('success', 'Siswa berhasil dihapus')
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menghapus data')
    }
  }

  return (
    <div>
      <div className="page-head"><h2>Data Siswa</h2></div>
      {msg && <div className={`alert alert-${msg.type}`}>{msg.text}</div>}

      <div className="grid-2">
        <div className="panel">
          <h3>{editing ? 'Ubah Siswa' : 'Tambah Siswa'}</h3>
          <form onSubmit={submit}>
            <label>NIS</label>
            <input
              value={form.nis}
              disabled={editing}
              onChange={(e) => setForm({ ...form, nis: e.target.value })}
              required
            />
            <label>Nama</label>
            <input
              value={form.nama}
              onChange={(e) => setForm({ ...form, nama: e.target.value })}
              required
            />
            <label>Kelas</label>
            <input
              value={form.kelas}
              onChange={(e) => setForm({ ...form, kelas: e.target.value })}
              required
            />
            <div className="form-actions">
              <button className="btn btn-primary">{editing ? 'Simpan' : 'Tambah'}</button>
              {editing && <button type="button" className="btn btn-ghost" onClick={cancel}>Batal</button>}
            </div>
          </form>
        </div>

        <div className="panel">
          <h3>Daftar Siswa ({list.length})</h3>
          {loading ? <p>Memuat...</p> : (
            <table className="table">
              <thead>
                <tr><th>NIS</th><th>Nama</th><th>Kelas</th><th>Aksi</th></tr>
              </thead>
              <tbody>
                {list.map((s) => (
                  <tr key={s.nis}>
                    <td>{s.nis}</td>
                    <td>{s.nama}</td>
                    <td>{s.kelas}</td>
                    <td className="row-actions">
                      <button className="btn btn-sm" onClick={() => edit(s)}>Ubah</button>
                      <button className="btn btn-sm btn-danger" onClick={() => remove(s.nis)}>Hapus</button>
                    </td>
                  </tr>
                ))}
                {list.length === 0 && <tr><td colSpan="4" className="empty">Belum ada data siswa</td></tr>}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  )
}
