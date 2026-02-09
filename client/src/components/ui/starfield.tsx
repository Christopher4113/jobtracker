
import { useEffect, useRef } from "react"

export function Starfield() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext("2d")
    if (!ctx) return

    let animationId: number

    const stars: { x: number; y: number; z: number; size: number; opacity: number }[] = []
    const STAR_COUNT = 400
    const SPEED = 0.15

    function resize() {
      if (!canvas) return
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }

    function initStars() {
      stars.length = 0
      for (let i = 0; i < STAR_COUNT; i++) {
        stars.push({
          x: Math.random() * (canvas?.width ?? 1920) - (canvas?.width ?? 1920) / 2,
          y: Math.random() * (canvas?.height ?? 1080) - (canvas?.height ?? 1080) / 2,
          z: Math.random() * 1000,
          size: Math.random() * 1.5 + 0.5,
          opacity: Math.random() * 0.5 + 0.5,
        })
      }
    }

    function animate() {
      if (!ctx || !canvas) return

      ctx.clearRect(0, 0, canvas.width, canvas.height)

      const cx = canvas.width / 2
      const cy = canvas.height / 2

      for (const star of stars) {
        star.z -= SPEED

        if (star.z <= 0) {
          star.x = Math.random() * canvas.width - cx
          star.y = Math.random() * canvas.height - cy
          star.z = 1000
        }

        const sx = (star.x / star.z) * 300 + cx
        const sy = (star.y / star.z) * 300 + cy

        if (sx < 0 || sx > canvas.width || sy < 0 || sy > canvas.height) {
          star.x = Math.random() * canvas.width - cx
          star.y = Math.random() * canvas.height - cy
          star.z = 1000
          continue
        }

        const brightness = 1 - star.z / 1000
        const size = star.size * (1 - star.z / 1000) * 2

        ctx.beginPath()
        ctx.arc(sx, sy, Math.max(size, 0.3), 0, Math.PI * 2)
        ctx.fillStyle = `rgba(200, 220, 255, ${brightness * star.opacity})`
        ctx.fill()
      }

      animationId = requestAnimationFrame(animate)
    }

    resize()
    initStars()
    animate()

    window.addEventListener("resize", () => {
      resize()
      initStars()
    })

    return () => {
      cancelAnimationFrame(animationId)
      window.removeEventListener("resize", resize)
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      className="fixed inset-0 -z-10"
      aria-hidden="true"
    />
  )
}
