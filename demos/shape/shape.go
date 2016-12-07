package shape

import (
	"encoding/binary"
	"math/rand"

	"math"

	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/drawable"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/geom"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util"
	"github.com/omustardo/gome/util/bytecoder"
)

type Shape interface {
	// Draw draws an outline of a Shape using line segments.
	DrawWireframe()
	// DrawFilled draws a filled in Shape using triangles.
	DrawFilled()
	// SetCenter sets the Shape to the specified position.
	SetCenter(x, y float32)
	// ModifyCenter moves the Shape by the specified amount.
	ModifyCenter(x, y float32)
	// Center is a point about which all actions, like rotation, are defined.
	// TODO: Consider the ability to modify the center point for rotating.
	Center() mgl32.Vec3
}

type ParallaxRect struct {
	geom.Rect
	Target camera.Camera
	// Essentially, how this object moves in comparison to the camera.
	// 1 is the same speed. 0.2 is 20% of camera speed.
	// The larger the number, the further away the object appears to be. For example, a ratio of 0.95 means the object
	// barely move when the camera moves - just like something that's very far away.
	// Negative numbers will make it move in the opposite direction, which isn't recommended.
	TranslationRatio float32
}

func GenParallaxRects(target camera.Camera, count int, minWidth, maxWidth, minSpeedRatio, maxSpeedRatio float32) []ParallaxRect {
	shapes := make([]ParallaxRect, count)
	for i := 0; i < count; i++ {
		shapes[i] = ParallaxRect{
			Rect: geom.Rect{
				Entity: entity.Entity{
					Position: mgl32.Vec3{rand.Float32()*20000 - 10000, rand.Float32()*20000 - 10000, 0},
					Scale:    mgl32.Vec3{rand.Float32()*(maxWidth-minWidth) + minWidth, rand.Float32()*(maxWidth-minWidth) + minWidth, 0},
					Rotation: mgl32.Vec3{0, 0, rand.Float32() * 2 * math.Pi},
				},
				Drawable: drawable.Drawable{
					Color: &color.RGBA{util.RandUint8(), util.RandUint8(), util.RandUint8(), 255},
				},
			},
			Target:           target,
			TranslationRatio: rand.Float32()*(maxSpeedRatio-minSpeedRatio) + minSpeedRatio,
		}
	}
	return shapes
}

func GetParallaxBuffers(arr []ParallaxRect) (parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer, parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer gl.Buffer) {
	lower, upper := float32(-0.5), float32(0.5)
	vertices := []float32{
		// Triangle 1
		lower, lower, 0,
		upper, upper, 0,
		lower, upper, 0,
		// Triangle 2
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
	}
	var posData, transData, transRatioData, angleData, scaleData, colorData []float32
	for _, rect := range arr {
		posData = append(posData, vertices...)
		for i := 0; i < 6; i++ {
			tx, ty, _ := rect.Center().Elem()
			transData = append(transData, tx, ty)
			transRatioData = append(transRatioData, rect.TranslationRatio)
			angleData = append(angleData, rect.Rotation.Z())
			scaleData = append(scaleData, rect.Scale.X(), rect.Scale.Y())
			colorData = append(colorData, float32(rect.Color.R)/255, float32(rect.Color.G)/255, float32(rect.Color.B)/255, float32(rect.Color.A)/255)
		}
	}
	parallaxPositionBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxPositionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, posData...), gl.STATIC_DRAW)

	parallaxTranslationBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxTranslationBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, transData...), gl.STATIC_DRAW)

	parallaxTranslationRatioBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxTranslationRatioBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, transRatioData...), gl.STATIC_DRAW)

	parallaxAngleBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxAngleBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, angleData...), gl.STATIC_DRAW)

	parallaxScaleBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxScaleBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, scaleData...), gl.STATIC_DRAW)

	parallaxColorBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxColorBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, colorData...), gl.STATIC_DRAW)
	return
}

func DrawParallaxBuffers(numObjects int, camPos mgl32.Vec3, parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer, parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer gl.Buffer) {
	gl.UseProgram(shader.Parallax.Program)
	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxPositionBuffer)
	gl.VertexAttribPointer(shader.Parallax.PositionAttrib, 3 /* floats per vertex */, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.PositionAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxTranslationBuffer)
	gl.VertexAttribPointer(shader.Parallax.TranslationAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.TranslationAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxTranslationRatioBuffer)
	gl.VertexAttribPointer(shader.Parallax.TranslationRatioAttrib, 1, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.TranslationRatioAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxAngleBuffer)
	gl.VertexAttribPointer(shader.Parallax.AngleAttrib, 1, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.AngleAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxScaleBuffer)
	gl.VertexAttribPointer(shader.Parallax.ScaleAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.ScaleAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, parallaxColorBuffer)
	gl.VertexAttribPointer(shader.Parallax.ColorAttrib, 4, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Parallax.ColorAttrib)

	shader.Parallax.SetCameraPosition(camPos)
	gl.DrawArrays(gl.TRIANGLES, 0, numObjects)
}

type OrbitingRect struct {
	geom.Rect
	// milliseconds to go entirely around the orbit. i.e. one year for the earth.
	// Goes counterclockwise by default. Set negative to go clockwise.
	revolutionSpeed int64
	orbit           geom.Circle
	// Makes the center of the orbit an object that can move. If nil, just uses the orbit's static center.
	orbitTarget entity.Target
	// rotateSpeed is the milliseconds to do one full rotation. i.e. one day for the earth.
	// Goes counterclockwise by default. Set negative to go clockwise. Use 0 to not rotate.
	rotateSpeed int64
}

func NewOrbitingRect(rect geom.Rect, orbitCenter mgl32.Vec2, orbitRadius float32, orbitTarget entity.Target, revolutionSpeed, rotateSpeed int64) *OrbitingRect {
	r := &OrbitingRect{
		Rect: rect,
		orbit: geom.Circle{
			Entity: entity.Entity{
				Position: orbitCenter.Vec3(0),
				Scale:    mgl32.Vec3{orbitRadius, orbitRadius, 0},
				Rotation: mgl32.Vec3{},
			},
			Drawable: drawable.Drawable{
				Color: &color.RGBA{140, 140, 140, 255},
			},
		},
		orbitTarget:     orbitTarget,
		revolutionSpeed: revolutionSpeed,
		rotateSpeed:     rotateSpeed,
	}
	r.Update()
	return r
}

func (r *OrbitingRect) Update() {
	if r.orbitTarget != nil {
		r.orbit.Position = r.orbitTarget.Center()
	}
	now := util.GetTimeMillis()
	percentRevolution := float32(now%r.revolutionSpeed) / float32(r.revolutionSpeed)
	rads := percentRevolution * 2 * math.Pi
	offset := mgl32.Vec3{float32(math.Cos(float64(rads))), float32(math.Sin(float64(rads))), 0}.Mul(r.orbit.Scale[0]) // TODO: use multiple scale dimensions (essentially the radius) to have elliptical orbits
	x, y, _ := r.orbit.Center().Add(offset).Elem()
	r.SetCenter(x, y, 0)

	// Elliptical orbit calculation: http://math.stackexchange.com/questions/22064/calculating-a-point-that-lies-on-an-ellipse-given-an-angle
	//a, b := float64(r.orbit.Scale[0]), float64(r.orbit.Scale[1])
	//x := a * b / math.Sqrt(b*b+a*a*math.Tan(rads)*math.Tan(rads))
	//y := a * b / math.Sqrt(a*a+b*b/math.Tan(rads)*math.Tan(rads))
	//if -math.Pi/2 <= rads || rads >= math.Pi/2 {
	//	x = -x
	//	y = -y
	//}

	if r.rotateSpeed != 0 {
		percentRotation := float32(now%r.rotateSpeed) / float32(r.rotateSpeed)
		r.Rotation[2] = percentRotation * 2 * math.Pi
	}
}

func (r *OrbitingRect) DrawOrbit() {
	r.orbit.DrawWireframe()
}
