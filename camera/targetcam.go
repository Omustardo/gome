package camera

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/entity"
)

var _ Camera = (*TargetCamera)(nil)

// TargetCamera is a camera that is always positioned at an offset from the target entity. It always looks toward that
// entity with the provided Up vector being up.
type TargetCamera struct {
	Target entity.Entity
	// targetOffset is where the camera is positioned in relation to the target.
	TargetOffset mgl32.Vec3
	Up           mgl32.Vec3
	Near, Far    float32
	Zoomer       zoom.Zoom

	// Field of view - only matters if using ProjectionPerspective.
	FOV float32
}

//
//func NewTargetCamera(target entity.Entity, targetOffset, up mgl32.Vec3, zoomer zoom.Zoom, near, far float32) (Camera, error) {
//	if near >= far {
//		return nil, fmt.Errorf("near(%v) >= far(%v)", near, far)
//	}
//	basicCam, err := NewBasicCamera(target.Center(), near, far)
//	if err != nil {
//		return nil, err
//	}
//	p := &TargetCamera{
//		BasicCamera: basicCam,
//		Target:      target,
//		zoomer:      zoomer,
//	}
//
//	p.Update()
//	return p, nil
//}

func (c *TargetCamera) ModelView() mgl32.Mat4 {
	targetPos := c.Target.Center()
	return mgl32.LookAt(
		targetPos.X(), targetPos.Y(), 1, // Camera Position. Always above target.
		targetPos.X(), targetPos.Y(), 0, // Target Position.
		0, 1, 0) // Up vector // TODO: For a camera that has Screen Up always as the direction the entity is facing, I think we just need to modify this line.
}

func (c *TargetCamera) ProjectionOrthographic(width, height float32) mgl32.Mat4 {
	zoomPercent := c.GetCurrentZoomPercent()
	width /= zoomPercent
	height /= zoomPercent
	return mgl32.Ortho(-width/2, width/2,
		-height/2, height/2,
		c.Near, c.Far)
}

func (c *TargetCamera) ProjectionPerspective(width, height float32) mgl32.Mat4 {
	return mgl32.Mat4{} // TODO
	//zoomPercent := c.GetCurrentZoomPercent()
	//width /= zoomPercent
	//height /= zoomPercent
	//if c.FOV <= 0 {
	//	log.Printf("invalid Field of View (%v). Using 45 degrees.", c.FOV)
	//	return mgl32.Perspective(45, float32(width)/float32(height), c.Near, c.Far)
	//}
	//return mgl32.Perspective(c.FOV, float32(width)/float32(height), c.Near, c.Far)
}

func (c *TargetCamera) Update() {
	if c.Zoomer != nil {
		c.Zoomer.Update()
	}
	if c.Near >= c.Far {
		log.Println("camera's near is closer than far - nothing will render")
	}
}

func (c *TargetCamera) Position() mgl32.Vec3 {
	return c.Target.Center().Add(c.TargetOffset)
}

// ScreenToWorldCoord2D returns the world coordinates of a point on the screen.
// This depends on the camera always looking directly down onto the XY-plane. e.g. camera position has positive Z.
// The screen space coordinate is expected in the coordinate system where the top left corner is (0,0), Y increases down, and X increases right.
func (c *TargetCamera) ScreenToWorldCoord2D(screenPoint mgl32.Vec2, windowWidth, windowHeight int) mgl32.Vec2 {
	zoomPercent := c.GetCurrentZoomPercent()
	pos := c.Position()
	return mgl32.Vec2{
		pos.X() + (screenPoint.X()-float32(windowWidth)/2)/zoomPercent,
		pos.Y() - (screenPoint.Y()-float32(windowHeight)/2)/zoomPercent,
	}
}

func (c *TargetCamera) GetCurrentZoomPercent() float32 {
	if c.Zoomer == nil {
		return 1
	}
	zoomPercent := c.Zoomer.GetCurrentPercent()
	if zoomPercent <= 0 {
		log.Printf("Invalid camera zoom: %v. Using default", zoomPercent)
		zoomPercent = 1.0
	}
	return zoomPercent
}
