const API_URL = import.meta.env.VITE_API_URL as string

export type AuthResponse = {
  token: string
  user: { id: string; name: string; email: string }
}

export async function signup(payload: { name: string; email: string; password: string }) {
  const res = await fetch(`${API_URL}/api/auth/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })

  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Signup failed")
  return data as AuthResponse
}

export async function login(payload: { email: string; password: string }) {
  const res = await fetch(`${API_URL}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })

  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Login failed")
  return data as AuthResponse
}

export function setToken(token: string) {
  sessionStorage.setItem("token", token)
}

export function getToken() {
  return sessionStorage.getItem("token")
}

export function clearToken() {
  sessionStorage.removeItem("token")
}

export async function me() {
  const token = getToken()
  const res = await fetch(`${API_URL}/api/me`, {
    headers: {
      Authorization: token ? `Bearer ${token}` : "",
    },
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Unauthorized")
  return data as { user: { id: string; name: string; email: string } }
}
