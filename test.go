package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
)

type Vec3 struct {
	x, y, z float64
}

type shape interface {
	test(r Vec3) float64
	color() Vec3
}

type sphere struct {
	p, t, c Vec3
	r       float64
}

func (s *sphere) test(r Vec3) float64 {
	rel := rotXYZ(r, s.t)
	t := Vec3{
		x: rel.x - s.p.x,
		y: rel.y - s.p.y,
		z: rel.z - s.p.z,
	}
	return length(t) - s.r
}

func (s *sphere) color() Vec3 {
	return s.c
}

type box struct {
	p, t, c Vec3
}

func (b *box) test(r Vec3) float64 {
	rel := rotXYZ(r, b.t)
	d := Vec3{
		x: math.Abs(rel.x) - b.p.x,
		y: math.Abs(rel.y) - b.p.y,
		z: math.Abs(rel.z) - b.p.z,
	}
	return length(maxVec3(d, 0.)) + math.Min(math.Max(d.x, math.Max(d.y, d.z)), 0.)
}

func (b *box) color() Vec3 {
	return b.c
}

type torus struct {
	p, t, c Vec3
	x, y    float64
}

func (t *torus) test(r Vec3) float64 {
	rel := rotXYZ(r, t.t)
	qx := length2(rel.x, rel.z) - t.x
	qy := rel.y
	return length2(qx, qy) - t.y
}

func (t *torus) color() Vec3 {
	return t.c
}

const SCALE int = 1
const HEIGHT = 800 * SCALE
const WIDTH = 800 * SCALE
const HF = float64(HEIGHT)
const WF = float64(WIDTH)

const MAXITER int = 32
const HIT float64 = 0.00001
const E float64 = 0.0005

func calcPos(ro, rd Vec3, dist float64) Vec3 {
	return Vec3{
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

func rotX(v Vec3, a float64) Vec3 {
	return Vec3{
		x: v.x,
		y: v.y*math.Cos(a) + v.z*-math.Sin(a),
		z: v.y*math.Sin(a) + v.z*math.Cos(a),
	}
}

func rotY(v Vec3, a float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(a) + v.z*math.Sin(a),
		y: v.y,
		z: v.x*-math.Sin(a) + v.z*math.Cos(a),
	}
}

func rotZ(v Vec3, a float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(a) + v.y*-math.Sin(a),
		y: v.x*math.Sin(a) + v.y*math.Cos(a),
		z: v.z,
	}
}

func rotXYZ(v, t Vec3) Vec3 {
	return rotX(rotY(rotZ(v, t.z), t.y), t.x)
}

func sphereTest(r Vec3) float64 {
	s := Vec3{x: 0.0, y: 0.0, z: 0.0}
	t := Vec3{
		x: r.x - s.x,
		y: r.y - s.y,
		z: r.z - s.z,
	}
	return length(t) - 1.5
}

func boxTest(r Vec3) float64 {
	b := Vec3{x: 1., y: 1., z: 0.5}
	d := Vec3{
		x: math.Abs(r.x) - b.x,
		y: math.Abs(r.y) - b.y,
		z: math.Abs(r.z) - b.z,
	}
	return length(maxVec3(d, 0.)) + math.Min(math.Max(d.x, math.Max(d.y, d.z)), 0.)
}

func maxVec3(v Vec3, m float64) Vec3 {
	return Vec3{
		x: math.Max(v.x, m),
		y: math.Max(v.y, m),
		z: math.Max(v.z, m),
	}
}

func length(v Vec3) float64 {
	return math.Sqrt(math.Pow(v.x, 2.0) + math.Pow(v.y, 2.0) + math.Pow(v.z, 2.0))
}
func length2(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2.0) + math.Pow(y, 2.0))
}

func normalize(v Vec3) Vec3 {
	l := length(v)
	return Vec3{x: v.x / l, y: v.y / l, z: v.z / l}
}

func dot(a, b Vec3) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func shapeColor(diff float64, col Vec3) color.RGBA {
	return color.RGBA{
		R: uint8(255.0 * col.x * diff),
		G: uint8(255.0 * col.y * diff),
		B: uint8(255.0 * col.z * diff),
		A: 255,
	}
}

func estimateNormal(r Vec3, sh shape) Vec3 {
	return normalize(Vec3{
		x: sh.test(Vec3{r.x + E, r.y, r.z}) - sh.test(Vec3{r.x - E, r.y, r.z}),
		y: sh.test(Vec3{r.x, r.y + E, r.z}) - sh.test(Vec3{r.x, r.y - E, r.z}),
		z: sh.test(Vec3{r.x, r.y, r.z + E}) - sh.test(Vec3{r.x, r.y, r.z - E}),
	})
}

func macroTest(r Vec3, shapes ...shape) (float64, Vec3) {
	res := 1000000000.0
	var col Vec3
	var val float64
	for _, sh := range shapes {
		val = math.Min(sh.test(r), res)
		if val < res {
			res = val
			col = sh.color()
		}
	}
	return res, col
}

func macroEstNorm(r Vec3, shapes ...shape) Vec3 {
	res := Vec3{x: 1000000000.0, y: 1000000000.0, z: 1000000000.0}
	for _, sh := range shapes {
		n := estimateNormal(r, sh)
		res.x = math.Min(n.x, res.x)
		res.y = math.Min(n.y, res.y)
		res.z = math.Min(n.z, res.z)
	}
	return res
}

func march(ro Vec3, rd Vec3) color.RGBA {
	lightPos := Vec3{x: 0.0, y: 0.0, z: -5.0}
	r, g, b := 0.0, 0.0, 0.0
	distTraveled := 0.0
	sh1 := torus{
		p: Vec3{x: 1.0, y: 1.0, z: 1.0},
		t: Vec3{x: 45.0, y: 45.0, z: 45.0},
		c: Vec3{x: 0.5, y: 0.25, z: 0.1},
		x: 3.5,
		y: 0.25,
	}
	sh2 := torus{
		p: Vec3{x: 2.0, y: 1.0, z: 1.0},
		t: Vec3{x: -90.0, y: 0.0, z: -90.0},
		c: Vec3{x: 0.1, y: 0.75, z: 0.25},
		x: 2.25,
		y: 0.15,
	}
	sh3 := sphere{
		p: Vec3{x: 0.0, y: 0.0, z: 0.0},
		t: Vec3{x: 0.0, y: 0.0, z: 0.0},
		c: Vec3{x: 0.1, y: 0.25, z: 0.75},
		r: 0.75,
	}
	for step := 0; step < MAXITER; step++ {
		finalPos := calcPos(ro, rd, distTraveled)
		dist, col := macroTest(finalPos, &sh1, &sh2, &sh3)
		if dist < HIT {
			normal := macroEstNorm(finalPos, &sh1, &sh2, &sh3)
			lightDir := normalize(Vec3{x: lightPos.x - finalPos.x, y: lightPos.y - finalPos.y, z: lightPos.z - finalPos.z})
			diff := math.Max(dot(normal, lightDir), 0.0)
			return shapeColor(diff, col)
		}
		distTraveled += dist
	}

	return color.RGBA{
		R: uint8(255.0 * r),
		G: uint8(255.0 * g),
		B: uint8(255.0 * b),
		A: 255,
	}
}

type packet struct {
	col  color.RGBA
	x, y int
}

func marchgo(chn chan packet) {
	go func() {
		for x := 0; x < WIDTH; x++ {
			xf := ((float64(x) / WF) * 2.0) - 1.0
			go func(x int, xf float64) {
				for y := 0; y < HEIGHT; y++ {
					ro := Vec3{x: 0.0, y: 0.0, z: -10.0}
					yf := ((float64(y) / HF) * 2.0) - 1.0
					rd := Vec3{x: xf, y: yf, z: 1.0}
					chn <- packet{col: march(ro, rd), x: x, y: y}
				}
			}(x, xf)
		}
	}()
}

func main() {
	file, err := os.Create("temp.png")
	check(err)
	img := image.NewNRGBA64(image.Rect(0, 0, WIDTH, HEIGHT))
	chn := make(chan packet)

	start := time.Now()
	go marchgo(chn)
	limit := WIDTH * HEIGHT
	for i := 0; i < limit; i++ {
		pack := <-chn
		img.Set(pack.x, pack.y, pack.col)
	}
	end := time.Now()

	fmt.Printf("took %v\n", end.Sub(start))
	err = png.Encode(file, img)
	check(err)
}
