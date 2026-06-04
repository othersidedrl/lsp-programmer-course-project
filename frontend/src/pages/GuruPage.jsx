import { useEffect, useState } from 'react'
import api from '../api/client'

const empty = { id_guru: '', nama: '', mata_pelajaran: '' }

export default function GuruPage() {
  const [list, setList] = useState([])
  const [form, setForm] = useState(empty)
  const [editing, setEditing] = useState(false)
  const [msg, setMsg] = useState(null)

  const notify = (type, text) => {
    setMsg({ type, text })
    setTimeout(() => setMsg(null), 3000)
  }

  const load = async () => {
    const res = await api.get('/guru')
    setList(res.data.data || [])
  }
  useEffect(() => { load() }, [])

  const submit = async (e) => {
    e.preventDefault()
    try {
      if (editing) {
        await api.put(`/guru/${form.id_guru}`, form)
        notify('success', 'Data guru berhasil diperbarui')
      } else {
        await api.post('/guru', form)
        notify('success', 'Guru berhasil ditambahkan')
      }
      setForm(empty)
      setEditing(false)
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menyimpan data')
    }
  }

  const edit = (g) => { setForm(g); setEditing(true) }
  const cancel = () => { setForm(empty); setEditing(false) }

  const remove = async (id) => {
    if (!window.confirm(`Hapus guru ${id}?`)) return
    try {
      await api.delete(`/guru/${id}`)
      notify('success', 'Guru berhasil dihapus')
      load()
    } catch (err) {
      notify('error', err.response?.data?.message || 'Gagal menghapus data')
    }
  }

  return (
    <div>
      <div className="page-head"><h2>Data Guru</h2></div>
      {msg && <div className={`alert alert-${msg.type}`}>{msg.text}</div>}

      <div className="grid-2">
        <div className="panel">
          <h3>{editing ? 'Ubah Guru' : 'Tambah Guru'}</h3>
          <form onSubmit={submit}>
            <label>ID Guru</label>
            <input
              value={form.id_guru}
              disabled={editing}
              onChange={(e) => setForm({ ...form, id_guru: e.target.value })}
              required
            />
            <label>Nama Guru</label>
            <input
              value={form.nama}
              onChange={(e) => setForm({ ...form, nama: e.target.value })}
              required
            />
            <label>Mata Pelajaran</label>
            <input
              value={form.mata_pelajaran}
              onChange={(e) => setForm({ ...form, mata_pelajaran: e.target.value })}
              required
            />
            <div className="form-actions">
              <button className="btn btn-primary">{editing ? 'Simpan' : 'Tambah'}</button>
              {editing && <button type="button" className="btn btn-ghost" onClick={cancel}>Batal</button>}
            </div>
          </form>
        </div>

        <div className="panel">
          <h3>Daftar Guru ({list.length})</h3>
          <table className="table">
            <thead>
              <tr><th>ID</th><th>Nama</th><th>Mata Pelajaran</th><th>Aksi</th></tr>
            </thead>
            <tbody>
              {list.map((g) => (
                <tr key={g.id_guru}>
                  <td>{g.id_guru}</td>
                  <td>{g.nama}</td>
                  <td>{g.mata_pelajaran}</td>
                  <td className="row-actions">
                    <button className="btn btn-sm" onClick={() => edit(g)}>Ubah</button>
                    <button className="btn btn-sm btn-danger" onClick={() => remove(g.id_guru)}>Hapus</button>
                  </td>
                </tr>
              ))}
              {list.length === 0 && <tr><td colSpan="4" className="empty">Belum ada data guru</td></tr>}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
