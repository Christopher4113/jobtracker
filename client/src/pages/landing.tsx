import { Starfield } from "@/components/ui/starfield"
import { Button } from "@/components/ui/button"
import { Rocket } from "lucide-react"
import { Link } from "react-router-dom";

const Landing = () => {
  return (
    <div className="relative min-h-screen overflow-hidden bg-gradient-to-b from-black via-[#050816] to-[#0b1020]">
      
      {/* Space glow */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(120,119,198,0.15),_transparent_70%)]" />

      {/* Stars */}
      <Starfield />

      <div className="relative z-10 flex min-h-screen flex-col items-center justify-center gap-8 px-4 text-center">
        <div className="flex items-center gap-3">
          <Rocket className="h-10 w-10 text-primary md:h-12 md:w-12 " />
          <h1 className="text-5xl font-bold tracking-tight text-foreground md:text-7xl text-balance text-white">
            Job Tracker
          </h1>
        </div>

        <p className="max-w-md text-lg text-muted-foreground">
          Navigate your career journey through the stars
        </p>

        <Link to="/login">
          <Button size="lg" className="mt-4 px-8 py-6 text-lg">
            Get Started
          </Button>
        </Link>
      </div>
    </div>
  )
}

export default Landing
