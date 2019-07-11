package main

import (
	"math"
)

/*
	Shapes
*/

type shape interface {
	test(r Vec3) float64
	color() Vec3
}

// SPHERE
type sphere struct {
	p, t, a, c Vec3
	r          float64
	preFuncs   []func(Vec3) Vec3
	postFuncs  []func(float64) float64
}

func (s *sphere) test(r Vec3) float64 {
	for _, fn := range s.preFuncs {
		r = fn(r)
	}
	rel := rotXYZ(affine(r, s.a), s.t)
	t := Vec3{
		x: rel.x - s.p.x,
		y: rel.y - s.p.y,
		z: rel.z - s.p.z,
	}
	v := length(t) - s.r
	for _, fn := range s.postFuncs {
		v = fn(v)
	}
	return v
}

func (s *sphere) color() Vec3 {
	return s.c
}

// BOX
type box struct {
	p, t, a, c Vec3
	preFuncs   []func(Vec3) Vec3
	postFuncs  []func(float64) float64
}

func (b *box) test(r Vec3) float64 {
	for _, fn := range b.preFuncs {
		r = fn(r)
	}
	rel := rotXYZ(affine(r, b.a), b.t)
	d := Vec3{
		x: math.Abs(rel.x) - b.p.x,
		y: math.Abs(rel.y) - b.p.y,
		z: math.Abs(rel.z) - b.p.z,
	}
	v := length(maxVec3(d, 0.)) + math.Min(math.Max(d.x, math.Max(d.y, d.z)), 0.)
	for _, fn := range b.postFuncs {
		v = fn(v)
	}
	return v
}

func (b *box) color() Vec3 {
	return b.c
}

// TORUS
type torus struct {
	p, t, a, c Vec3
	x, y       float64
	preFuncs   []func(Vec3) Vec3
	postFuncs  []func(float64) float64
}

func (t *torus) test(r Vec3) float64 {
	for _, fn := range t.preFuncs {
		r = fn(r)
	}
	rel := rotXYZ(affine(r, t.a), t.t)
	qx := length2(rel.x, rel.z) - t.x
	qy := rel.y
	v := length2(qx, qy) - t.y
	for _, fn := range t.postFuncs {
		v = fn(v)
	}
	return v
}

func (t *torus) color() Vec3 {
	return t.c
}

//SERPENSKI
type serpenski struct {
	p, t, a, c Vec3
}

//func (s *serpenski) test(r Vec3) float64 {
//	var v float64
//	iters, n := 10, 0
//	for n < iters; n++ {
//		if (r.x+r.y) < 0. {
//			r.x, r.y = -r.y, -r.x
//		}
//		if (r.x+r.z) < 0. {
//			r.x, r.z = -r.z, -r.x
//		}
//		if (r.y+r.z) < 0. {
//			r.y, r.z = -r.z, -r.y
//		}
//		r = r*scale -
//	}
//
//}
func (s *serpenski) test(r Vec3) float64 {
	z := r
	size := 1.0
	a1 := Vec3{x: size, y: size, z: size}
	a2 := Vec3{x: -size, y: -size, z: size}
	a3 := Vec3{x: size, y: -size, z: -size}
	a4 := Vec3{x: -size, y: size, z: -size}
	var c Vec3
	n := 0
	dist := 0.0
	d := 0.0
	scale := 2.0
	scm1 := scale - 1.0
	iters := 10
	for n < iters {
		c = a1
		dist = length(Vec3{x: z.x - a1.x, y: z.y - a1.y, z: z.z - a1.z})
		d = length(Vec3{x: z.x - a2.x, y: z.y - a2.y, z: z.z - a2.z})
		if d < dist {
			c = a2
			dist = d
		}
		d = length(Vec3{x: z.x - a3.x, y: z.y - a3.y, z: z.z - a3.z})
		if d < dist {
			c = a3
			dist = d
		}
		d = length(Vec3{x: z.x - a4.x, y: z.y - a4.y, z: z.z - a4.z})
		if d < dist {
			c = a4
			dist = d
		}
		z = Vec3{x: (z.x * scale) - (c.x * scm1), y: (z.y * scale) - (c.y * scm1), z: (z.z * scale) - (c.z * scm1)}
		n++
	}
	re := length(z) * math.Pow(scale, float64(-n))
	//fmt.Printf("%v\n", re)
	return re
}

func (s *serpenski) color() Vec3 {
	return s.c
}
