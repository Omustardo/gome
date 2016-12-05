package geom

// Loads models into buffers on the GPU. Must be called after glfw.Init()
func Initialize() {
	initializeCircle()
	initializeCube()
	initializeLine()
	initializeRect()
	initializeTriangle()
}
