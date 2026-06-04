import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

// ProtectedRoute melindungi rute: harus login, dan (opsional) sesuai peran.
export default function ProtectedRoute({ roles, children }) {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" replace />
  }
  if (roles && !roles.includes(user.role)) {
    return <Navigate to="/" replace />
  }
  return children
}
