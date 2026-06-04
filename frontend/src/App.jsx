import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'
import Layout from './components/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import SiswaPage from './pages/SiswaPage'
import GuruPage from './pages/GuruPage'
import NilaiPage from './pages/NilaiPage'
import LaporanPage from './pages/LaporanPage'
import NilaiSayaPage from './pages/NilaiSayaPage'

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />

          {/* Semua rute di bawah ini wajib login */}
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Dashboard />} />

            {/* Admin saja */}
            <Route path="siswa" element={<ProtectedRoute roles={['admin']}><SiswaPage /></ProtectedRoute>} />
            <Route path="guru" element={<ProtectedRoute roles={['admin']}><GuruPage /></ProtectedRoute>} />

            {/* Admin & Guru */}
            <Route path="nilai" element={<ProtectedRoute roles={['admin', 'guru']}><NilaiPage /></ProtectedRoute>} />
            <Route path="laporan" element={<ProtectedRoute roles={['admin', 'guru']}><LaporanPage /></ProtectedRoute>} />

            {/* Siswa saja */}
            <Route path="nilai-saya" element={<ProtectedRoute roles={['siswa']}><NilaiSayaPage /></ProtectedRoute>} />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}
