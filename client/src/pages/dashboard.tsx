import { useEffect, useMemo, useState } from "react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  clearToken,
  getToken,
  createJob,
  deleteJob,
  getJobs,
  updateJob,
  type Job,
  type JobStatus,
} from "@/lib/api"

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  if (Number.isNaN(d.getTime())) return "-"
  return d.toLocaleDateString(undefined, { month: "short", day: "2-digit", year: "numeric" })
}

const statusOptions: { value: JobStatus; label: string }[] = [
  { value: "applied", label: "Applied" },
  { value: "interviewing", label: "Interviewing" },
  { value: "offer", label: "Offer" },
  { value: "rejected", label: "Rejected" },
]

function getErrorMessage(err: unknown): string {
  if (err instanceof Error) return err.message
  return "Something went wrong"
}

export default function Dashboard() {
  const navigate = useNavigate()

  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string>("")

  // form state
  const [editingId, setEditingId] = useState<string | null>(null)
  const [company, setCompany] = useState<string>("")
  const [role, setRole] = useState<string>("")
  const [location, setLocation] = useState<string>("")
  const [status, setStatus] = useState<JobStatus>("applied")
  const [link, setLink] = useState<string>("")
  const [source, setSource] = useState<string>("")
  const [notes, setNotes] = useState<string>("")
  const [saving, setSaving] = useState<boolean>(false)

  const sortedJobs = useMemo(() => {
    const list = Array.isArray(jobs) ? jobs : []
    return [...list].sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
  }, [jobs])

  async function refreshJobs() {
    setError("")
    const data = await getJobs()
    const list = Array.isArray(data.jobs) ? data.jobs : []
    setJobs(list)
  }

  useEffect(() => {
    if (!getToken()) {
      navigate("/login")
      return
    }

    ;(async () => {
      try {
        setLoading(true)
        await refreshJobs()
      } catch (e: unknown) {
        setError(getErrorMessage(e))
      } finally {
        setLoading(false)
      }
    })()
  }, [navigate])

  function resetForm() {
    setEditingId(null)
    setCompany("")
    setRole("")
    setLocation("")
    setStatus("applied")
    setLink("")
    setSource("")
    setNotes("")
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setError("")

    if (!company.trim() || !role.trim()) {
      setError("Company and role are required")
      return
    }

    try {
      setSaving(true)

      if (editingId) {
        await updateJob(editingId, {
          company: company.trim(),
          role: role.trim(),
          location: location.trim(),
          status,
          link: link.trim(),
          source: source.trim(),
          notes: notes.trim(),
        })
      } else {
        await createJob({
          company: company.trim(),
          role: role.trim(),
          location: location.trim(),
          status,
          link: link.trim() || undefined,
          source: source.trim() || undefined,
          notes: notes.trim() || undefined,
        })
      }

      await refreshJobs()
      resetForm()
    } catch (e: unknown) {
      setError(getErrorMessage(e))
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: string) {
    const ok = window.confirm("Delete this job?")
    if (!ok) return

    try {
      setError("")
      await deleteJob(id)
      await refreshJobs()
      if (editingId === id) resetForm()
    } catch (e: unknown) {
      setError(getErrorMessage(e))
    }
  }

  function startEdit(j: Job) {
    setEditingId(j.id)
    setCompany(j.company ?? "")
    setRole(j.role ?? "")
    setLocation(j.location ?? "")
    setStatus(j.status)
    setLink(j.link ?? "")
    setSource(j.source ?? "")
    setNotes(j.notes ?? "")
  }

  function logout() {
    clearToken()
    navigate("/login")
  }

  return (
    <div className="relative min-h-screen overflow-hidden bg-gradient-to-b from-black via-[#050816] to-[#0b1020]">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(120,119,198,0.18),_transparent_70%)]" />

      <div className="relative z-10 mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-6 px-4 py-10">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">Dashboard</h1>
            <p className="text-sm text-white/60">Track applications, status changes, and notes</p>
          </div>

          <div className="flex gap-2">
            <Button
              type="button"
              className="bg-white/10 text-white hover:bg-white/15"
              onClick={() => refreshJobs()}
              disabled={loading}
            >
              Refresh
            </Button>
            <Button type="button" className="bg-white/90 text-black hover:bg-white" onClick={logout}>
              Log out
            </Button>
          </div>
        </div>

        <div className="rounded-2xl border border-white/15 bg-black/60 p-6 shadow-2xl backdrop-blur">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">{editingId ? "Edit job" : "Add a job"}</h2>
            {editingId && (
              <Button
                type="button"
                className="bg-white/10 text-white hover:bg-white/15"
                onClick={resetForm}
              >
                Cancel edit
              </Button>
            )}
          </div>

          <form className="grid gap-3 md:grid-cols-2" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Company</label>
              <Input
                value={company}
                onChange={(e) => setCompany(e.target.value)}
                placeholder="Stripe"
                className="text-white placeholder:text-white/40"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Role</label>
              <Input
                value={role}
                onChange={(e) => setRole(e.target.value)}
                placeholder="Software Engineer"
                className="text-white placeholder:text-white/40"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Location</label>
              <Input
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="San Francisco, CA"
                className="text-white placeholder:text-white/40"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Process</label>
              <select
                value={status}
                onChange={(e) => setStatus(e.target.value as JobStatus)}
                className="h-10 w-full rounded-md border border-white/15 bg-black/40 px-3 text-white outline-none"
              >
                {statusOptions.map((s) => (
                  <option key={s.value} value={s.value} className="bg-black">
                    {s.label}
                  </option>
                ))}
              </select>
              <p className="text-xs text-white/50">Changing status updates the status date automatically</p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Link</label>
              <Input
                value={link}
                onChange={(e) => setLink(e.target.value)}
                placeholder="Job posting URL"
                className="text-white placeholder:text-white/40"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Source</label>
              <Input
                value={source}
                onChange={(e) => setSource(e.target.value)}
                placeholder="Referral, LinkedIn, YC jobs"
                className="text-white placeholder:text-white/40"
              />
            </div>

            <div className="space-y-2 md:col-span-2">
              <label className="text-sm font-medium text-white">Notes</label>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Recruiter name, next steps, comp notes, reminders..."
                className="min-h-[90px] w-full rounded-md border border-white/15 bg-black/40 px-3 py-2 text-white placeholder:text-white/40 outline-none"
              />
            </div>

            {error && (
              <div className="md:col-span-2">
                <p className="text-sm text-red-400">{error}</p>
              </div>
            )}

            <div className="md:col-span-2">
              <Button
                type="submit"
                size="lg"
                disabled={saving}
                className="w-full bg-white/90 text-black hover:bg-white disabled:opacity-60"
              >
                {saving ? "Saving..." : editingId ? "Save changes" : "Add job"}
              </Button>
            </div>
          </form>
        </div>

        <div className="rounded-2xl border border-white/15 bg-black/60 p-6 shadow-2xl backdrop-blur">
          <h2 className="mb-4 text-lg font-semibold text-white">Your jobs</h2>

          {loading ? (
            <p className="text-white/60">Loading...</p>
          ) : sortedJobs.length === 0 ? (
            <p className="text-white/60">No jobs yet. Add your first one above.</p>
          ) : (
            <div className="w-full overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="text-white/70">
                  <tr className="border-b border-white/10">
                    <th className="py-3 pr-4">Company</th>
                    <th className="py-3 pr-4">Role</th>
                    <th className="py-3 pr-4">Location</th>
                    <th className="py-3 pr-4">Process</th>
                    <th className="py-3 pr-4">Status date</th>
                    <th className="py-3 pr-4">Last updated</th>
                    <th className="py-3 pr-4">Actions</th>
                  </tr>
                </thead>
                <tbody className="text-white/90">
                  {sortedJobs.map((j) => (
                    <tr key={j.id} className="border-b border-white/10">
                      <td className="py-3 pr-4 font-medium">{j.company}</td>
                      <td className="py-3 pr-4">{j.role}</td>
                      <td className="py-3 pr-4">{j.location || "-"}</td>
                      <td className="py-3 pr-4 capitalize">{j.status}</td>
                      <td className="py-3 pr-4">{formatDate(j.statusUpdatedAt)}</td>
                      <td className="py-3 pr-4">{formatDate(j.updatedAt)}</td>
                      <td className="py-3 pr-4">
                        <div className="flex gap-2">
                          <Button
                            type="button"
                            className="bg-white/10 text-white hover:bg-white/15"
                            onClick={() => startEdit(j)}
                          >
                            Edit
                          </Button>
                          <Button
                            type="button"
                            className="bg-white/10 text-white hover:bg-white/15"
                            onClick={() => handleDelete(j.id)}
                          >
                            Delete
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        <p className="text-xs text-white/40">
          Status date updates whenever you change Applied, Interviewing, Offer, or Rejected.
        </p>
      </div>
    </div>
  )
}
