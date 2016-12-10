package geom

// Loads models into buffers on the GPU. glfw.Init() must be called before calling this.
func Initialize() {
	initializeCircle()
	initializeCube()
	initializeLine()
	initializeRect()
	initializeTriangle()
}
