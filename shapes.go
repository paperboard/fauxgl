package fauxgl

import "math"

func NewPlane() *Mesh {
	v1 := Vector{-0.5, -0.5, 0}
	v2 := Vector{0.5, -0.5, 0}
	v3 := Vector{0.5, 0.5, 0}
	v4 := Vector{-0.5, 0.5, 0}
	return NewTriangleMesh([]*Triangle{
		NewTriangleForPoints(v1, v2, v3),
		NewTriangleForPoints(v1, v3, v4),
	})
}

func NewCube() *Mesh {
	v := []Vector{
		{-1, -1, -1}, {-1, -1, 1}, {-1, 1, -1}, {-1, 1, 1},
		{1, -1, -1}, {1, -1, 1}, {1, 1, -1}, {1, 1, 1},
	}
	mesh := NewTriangleMesh([]*Triangle{
		NewTriangleForPoints(v[3], v[5], v[7]),
		NewTriangleForPoints(v[5], v[3], v[1]),
		NewTriangleForPoints(v[0], v[6], v[4]),
		NewTriangleForPoints(v[6], v[0], v[2]),
		NewTriangleForPoints(v[0], v[5], v[1]),
		NewTriangleForPoints(v[5], v[0], v[4]),
		NewTriangleForPoints(v[5], v[6], v[7]),
		NewTriangleForPoints(v[6], v[5], v[4]),
		NewTriangleForPoints(v[6], v[3], v[7]),
		NewTriangleForPoints(v[3], v[6], v[2]),
		NewTriangleForPoints(v[0], v[3], v[2]),
		NewTriangleForPoints(v[3], v[0], v[1]),
	})
	mesh.Transform(Scale(Vector{0.5, 0.5, 0.5}))
	return mesh
}

func NewCubeForBox(box Box) *Mesh {
	m := Translate(Vector{0.5, 0.5, 0.5})
	m = m.Scale(box.Size())
	m = m.Translate(box.Min)
	cube := NewCube()
	cube.Transform(m)
	return cube
}

func NewCubeOutlineForBox(box Box) *Mesh {
	x0 := box.Min.X
	y0 := box.Min.Y
	z0 := box.Min.Z
	x1 := box.Max.X
	y1 := box.Max.Y
	z1 := box.Max.Z
	return NewLineMesh([]*Line{
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y0, z1}),
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x0, y1, z1}),
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y0, z1}),
		NewLineForPoints(Vector{x1, y1, z0}, Vector{x1, y1, z1}),
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y1, z0}),
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x0, y1, z1}),
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y1, z0}),
		NewLineForPoints(Vector{x1, y0, z1}, Vector{x1, y1, z1}),
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x1, y0, z0}),
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x1, y1, z0}),
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x1, y0, z1}),
		NewLineForPoints(Vector{x0, y1, z1}, Vector{x1, y1, z1}),
	})
}

func NewSphere(lngStep, latStep int) *Mesh {
	var triangles []*Triangle
	for lat0 := -90; lat0 < 90; lat0 += latStep {
		lat1 := lat0 + latStep
		v0 := float64(lat0+90) / 180
		v1 := float64(lat1+90) / 180
		for lng0 := -180; lng0 < 180; lng0 += lngStep {
			lng1 := lng0 + lngStep
			u0 := float64(lng0+180) / 360
			u1 := float64(lng1+180) / 360
			if lng1 >= 180 {
				lng1 -= 360
			}
			p00 := LatLngToXYZ(float64(lat0), float64(lng0))
			p01 := LatLngToXYZ(float64(lat0), float64(lng1))
			p10 := LatLngToXYZ(float64(lat1), float64(lng0))
			p11 := LatLngToXYZ(float64(lat1), float64(lng1))
			t1 := NewTriangleForPoints(p00, p01, p11)
			t2 := NewTriangleForPoints(p00, p11, p10)
			if lat0 != -90 {
				t1.V1.Texture = Vector{u0, v0, 0}
				t1.V2.Texture = Vector{u1, v0, 0}
				t1.V3.Texture = Vector{u1, v1, 0}
				triangles = append(triangles, t1)
			}
			if lat1 != 90 {
				t2.V1.Texture = Vector{u0, v0, 0}
				t2.V2.Texture = Vector{u1, v1, 0}
				t2.V3.Texture = Vector{u0, v1, 0}
				triangles = append(triangles, t2)
			}
		}
	}
	return NewTriangleMesh(triangles)
}

func NewSphere2(detail int) *Mesh {
	var triangles []*Triangle
	v := []Vector{
		{0, 0, -1}, {1, 0, 0}, {0, -1, 0},
		{-1, 0, 0}, {0, 1, 0}, {0, 0, 1},
	}
	triangles = append(triangles, newSphereHelper(detail, v[0], v[3], v[4])...)
	triangles = append(triangles, newSphereHelper(detail, v[0], v[4], v[1])...)
	triangles = append(triangles, newSphereHelper(detail, v[5], v[4], v[3])...)
	triangles = append(triangles, newSphereHelper(detail, v[5], v[1], v[4])...)
	triangles = append(triangles, newSphereHelper(detail, v[2], v[3], v[0])...)
	triangles = append(triangles, newSphereHelper(detail, v[1], v[2], v[0])...)
	triangles = append(triangles, newSphereHelper(detail, v[3], v[2], v[5])...)
	triangles = append(triangles, newSphereHelper(detail, v[2], v[1], v[5])...)
	return NewTriangleMesh(triangles)
}

func newSphereHelper(detail int, v1, v2, v3 Vector) []*Triangle {
	if detail == 0 {
		t := NewTriangleForPoints(v1, v2, v3)
		return []*Triangle{t}
	}
	var triangles []*Triangle
	v12 := v1.Add(v2).DivScalar(2).Normalize()
	v13 := v1.Add(v3).DivScalar(2).Normalize()
	v23 := v2.Add(v3).DivScalar(2).Normalize()
	triangles = append(triangles, newSphereHelper(detail-1, v1, v12, v13)...)
	triangles = append(triangles, newSphereHelper(detail-1, v2, v23, v12)...)
	triangles = append(triangles, newSphereHelper(detail-1, v3, v13, v23)...)
	triangles = append(triangles, newSphereHelper(detail-1, v12, v23, v13)...)
	return triangles
}

func NewCylinder(step int, capped bool) *Mesh {
	var triangles []*Triangle
	for a0 := 0; a0 < 360; a0 += step {
		a1 := (a0 + step) % 360
		r0 := Radians(float64(a0))
		r1 := Radians(float64(a1))
		x0 := math.Cos(r0)
		y0 := math.Sin(r0)
		x1 := math.Cos(r1)
		y1 := math.Sin(r1)
		p00 := Vector{x0, y0, -0.5}
		p10 := Vector{x1, y1, -0.5}
		p11 := Vector{x1, y1, 0.5}
		p01 := Vector{x0, y0, 0.5}
		t1 := NewTriangleForPoints(p00, p10, p11)
		t2 := NewTriangleForPoints(p00, p11, p01)
		triangles = append(triangles, t1)
		triangles = append(triangles, t2)
		if capped {
			p0 := Vector{0, 0, -0.5}
			p1 := Vector{0, 0, 0.5}
			t3 := NewTriangleForPoints(p0, p10, p00)
			t4 := NewTriangleForPoints(p1, p01, p11)
			triangles = append(triangles, t3)
			triangles = append(triangles, t4)
		}
	}
	return NewTriangleMesh(triangles)
}

func NewCone(step int, capped bool) *Mesh {
	var triangles []*Triangle
	for a0 := 0; a0 < 360; a0 += step {
		a1 := (a0 + step) % 360
		r0 := Radians(float64(a0))
		r1 := Radians(float64(a1))
		x0 := math.Cos(r0)
		y0 := math.Sin(r0)
		x1 := math.Cos(r1)
		y1 := math.Sin(r1)
		p00 := Vector{x0, y0, -0.5}
		p10 := Vector{x1, y1, -0.5}
		p1 := Vector{0, 0, 0.5}
		t1 := NewTriangleForPoints(p00, p10, p1)
		triangles = append(triangles, t1)
		if capped {
			p0 := Vector{0, 0, -0.5}
			t2 := NewTriangleForPoints(p0, p10, p00)
			triangles = append(triangles, t2)
		}
	}
	return NewTriangleMesh(triangles)
}
