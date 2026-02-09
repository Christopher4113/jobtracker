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

export type JobStatus = "applied" | "interviewing" | "offer" | "rejected"

export type Job = {
  id: string
  userId: string
  company: string
  role: string
  location: string
  status: JobStatus
  statusUpdatedAt: string
  link?: string
  notes?: string
  source?: string
  createdAt: string
  updatedAt: string
}

function authHeaders() {
  const token = getToken()
  return {
    "Content-Type": "application/json",
    Authorization: token ? `Bearer ${token}` : "",
  }
}

export async function getJobs() {
  const res = await fetch(`${API_URL}/api/jobs`, { headers: authHeaders() })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Failed to load jobs")
  return data as { jobs: Job[] }
}

export async function createJob(payload: {
  company: string
  role: string
  location: string
  status: JobStatus
  link?: string
  notes?: string
  source?: string
}) {
  const res = await fetch(`${API_URL}/api/jobs`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(payload),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Failed to create job")
  return data as { job: Job }
}

export async function updateJob(id: string, payload: Partial<Omit<Job, "id" | "userId" | "createdAt" | "updatedAt">>) {
  const res = await fetch(`${API_URL}/api/jobs/${id}`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify(payload),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Failed to update job")
  return data as { ok: boolean }
}

export async function deleteJob(id: string) {
  const res = await fetch(`${API_URL}/api/jobs/${id}`, {
    method: "DELETE",
    headers: authHeaders(),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error ?? "Failed to delete job")
  return data as { ok: boolean }
}
