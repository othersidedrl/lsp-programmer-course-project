import { createContext, useContext, useState } from 'react'
import api from '../api/client'

const AuthContext = createContext(null)

// AuthProvider menyimpan status login & menyediakan fungsi login/logout.
export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const raw = localStorage.getItem('user')
    return raw ? JSON.parse(raw) : null
  })

  const login = async (username, password) => {
    const res = await api.post('/login', { username, password })
    const { token, user: profile } = res.data.data
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(profile))
    setUser(profile)
    return profile
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
