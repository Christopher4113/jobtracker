import { useState } from "react"
import { useNavigate, Link } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { login, setToken } from "@/lib/api"

function getErrorMessage(err: unknown): string {
  if (err instanceof Error) return err.message
  return "Something went wrong"
}

export default function Login() {
  const navigate = useNavigate()

  const [email, setEmail] = useState<string>("")
  const [password, setPassword] = useState<string>("")
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string>("")

  const handleLogin = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError("")

    if (!email.trim() || !password) {
      setError("Please enter email and password")
      return
    }

    try {
      setLoading(true)
      const data = await login({ email: email.trim(), password })
      setToken(data.token)
      navigate("/dashboard")
    } catch (err: unknown) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative min-h-screen overflow-hidden bg-gradient-to-b from-black via-[#050816] to-[#0b1020]">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(120,119,198,0.18),_transparent_70%)]" />

      <div className="relative z-10 flex min-h-screen items-center justify-center px-4">
        <div className="w-full max-w-md rounded-2xl border border-white/15 bg-black/60 p-6 shadow-2xl backdrop-blur">
          <div className="mb-6 text-center">
            <h1 className="text-3xl font-bold text-white">Welcome back</h1>
            <p className="mt-2 text-sm text-white/60">
              Log in to continue tracking your applications
            </p>
          </div>

          <form className="space-y-4" onSubmit={handleLogin}>
            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Email</label>
              <Input
                name="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                className="text-white placeholder:text-white/40"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-white">Password</label>
              <Input
                name="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="text-white placeholder:text-white/40"
                required
              />
            </div>

            {error && <p className="text-sm text-red-400 text-center">{error}</p>}

            <Button
              type="submit"
              size="lg"
              disabled={loading}
              className="w-full bg-white/90 text-black hover:bg-white disabled:opacity-60"
            >
              {loading ? "Logging in..." : "Log in"}
            </Button>

            <p className="text-center text-sm text-white/60">
              No account?{" "}
              <Link to="/signup" className="text-white hover:underline">
                Sign up
              </Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}
