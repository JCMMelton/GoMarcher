package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Vec3 struct {
	x, y, z float64
}

const HEIGHT int = 800
const WIDTH  int = 800
const HF = float64(HEIGHT)
const WF = float64(WIDTH)

const MAXITER int = 32
const HIT float64 = 0.00001
const E float64 = 0.0005

func calcPos(ro, rd Vec3, dist float64) Vec3 {
	return Vec3 {
		x: ro.x + (dist * rd.x),
		y: ro.y + (dist * rd.y),
		z: ro.z + (dist * rd.z),
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sphereTest(r Vec3) float64 {
	s := Vec3 {x: 0.0, y: 0.0, z: 0.0}
	t := Vec3{
		x: r.x - s.x,
		y: r.y - s.y,
		z: r.z - s.z,
	}
	return length(t)-1.5
}

func length(v Vec3) float64 {
	return math.Sqrt(math.Pow(v.x, 2.0)+math.Pow(v.y, 2.0)+math.Pow(v.z, 2.0))
}

func normalize(v Vec3) Vec3 {
	l := length(v)
	return Vec3{x: v.x/l, y: v.y/l, z: v.z/l}
}

func dot(a, b Vec3) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func estimateNormal(r Vec3) Vec3 {
	return normalize(Vec3 {
		x: sphereTest(Vec3{r.x+E, r.y, r.z})-sphereTest(Vec3{r.x-E, r.y, r.z}),
		y: sphereTest(Vec3{r.x, r.y+E, r.z})-sphereTest(Vec3{r.x, r.y-E, r.z}),
		z: sphereTest(Vec3{r.x, r.y, r.z+E})-sphereTest(Vec3{r.x, r.y, r.z-E}),
	})
}

func march(ro Vec3, rd Vec3) color.RGBA {
	lightPos := Vec3{x: 10.0, y: 100.0, z: -10.0}
	r, g, b := 0.0, 0.0, 0.0
	distTraveled := 0.0
	for step:=0; step<MAXITER; step++ {
		finalPos := calcPos(ro, rd, distTraveled)
		dist := sphereTest(finalPos)
		if dist < HIT {
			r = 0.1
			g = 0.7
			b = 0.5
			normal := estimateNormal(finalPos)
			lightDir := normalize(Vec3{x: lightPos.x-finalPos.x, y: lightPos.y-finalPos.y, z: lightPos.z-finalPos.z })
			diff := math.Max(dot(normal, lightDir), 0.0)
			return color.RGBA{
				R: uint8(255.0*r*diff),
				G: uint8(255.0*g*diff),
				B: uint8(255.0*b*diff),
				A: 255,
			}

		distTraveled += dist
	}

	return color.RGBA{
		R: uint8(255.0*r),
		G: uint8(255.0*g),
		B: uint8(255.0*b),
		A: 255,
	}
}

func main() {
	file, err := os.Create("temp.png")
	check(err)
	img := image.NewNRGBA64(image.Rect(0, 0, WIDTH, HEIGHT))
	for x:=0; x<WIDTH; x++ {
		xf := ((float64(x)/WF)*2.0)-1.0
		for y:=0; y<HEIGHT; y++ {
			ro  := Vec3{x: 0.0, y: 0.0, z: -10.0}
			yf  := ((float64(y)/HF)*2.0)-1.0
			rd  := Vec3{x: xf, y: yf, z: 1.0}
			col := march(ro, rd)
			img.Set(x, y, col)
		}
	}
	err = png.Encode(file, img)
	check(err)
}