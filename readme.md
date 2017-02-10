## Gome (Golang game engine)

### Overview

Gome is a pet project of mine. I took a graphics course in college and am now re-learning by implementing it myself.
This project is certainly pre-alpha so I don't recommend depending on anything in it.

I'm not attempting to develop a full featured engine like Unreal or Unity, but I hope to eventually support enough of
the basics to implement a networked game with simple physics. I chose Go because I appreciate its simplicity and find
it pleasant to work with. It also has GopherJS, which allows programs developed in Go to be compiled not only
into regular binaries, but also into javascript. With the `github.com/goxjs/gl` and `glfw` packages, this allows the same 
OpenGL shaders to be compiled for both desktop and webgl (see the gome shader package).

Some of my GopherJS demos are hosted at https://omustardo.github.io/ but note that not all of them were made using gome.

### Usage

I recommend looking in the demos folder. To run a demo:
```
go get github.com/omustardo/gome
```
Then for GopherJS:
```
gopherjs serve
# In your browser, navigate to http://localhost:8080/github.com/omustardo/gome/demos
# Open any demo subfolder
```
For Desktop:
```
cd <path/to/workspace>/github.com/omustardo/gome/demos/<demo>
go build
./<demo_name.exe> --base_dir=.
```
To use `go run` you need to provide the `base_dir` flag with the full path to the assets folder.
This is due to an ugly hack to load local files while maintaining the ability to load files via the same relative
paths on the web. I explain the issue in more depth in https://github.com/Omustardo/gome/blob/master/asset/asset.go

If anyone has suggestions to improve file loading, I'm very open to them.

### Glossary
 * Mesh: a bunch of 3d points are read in and stored in buffers on the GPU, and references to these buffers are put
 into mesh objects. Meshes can be loaded from file by the asset package, or created dynamically like in the mesh package.

 * Entity: information about an object in the world. Position, Rotation, Scale. 

 * Model: The combination of a Mesh and Entity. Uses the world information from the entity, and the graphics information
 from the mesh to render.

 * NRGBA: non-premultiplied RGBA color. In a premultiplied RGBA color, the RGB values are always between 0 and the 
 Alpha value, since they've been scaled down. For example, if your NRGBA color is red and a bit transparent 
 `[200,0,0,100]`  then the premultiplied RGBA version is `[200*(100/255),0,0,100] = [78,0,0,100]`. 
 I believe users will find NRGBA more natural to work with so that's the default. 
 Note that you should access the color values directly, i.e. `color.R`, rather than using the `RGBA()` function.
 Also note that RGBA commonly used without multiplying the alpha value, but I chose NRGBA so it's explicit.
 
### Gotcha's
 * The scale of models in the game is up to the user to set via entity.Entity's Scale field. Note that leaving it empty
 will mean your model has no physical size and so won't be shown. If you want your model to not be rendered, it's 
 more effective to set model.Hidden = true

 * Be careful when using mgl32. Its functions tend to return new structs rather than modifying the caller.
 For example, you might be tempted to modify an Entity's position using `e.Position.Add(mgl32.Vec3{1,1,0})` 
 but this doesn't modify `e.Position`. It just returns a new vector. 
 Instead, use `e.Position = e.Position.Add(mgl32.Vec3{1,1,0})` or `e.ModifyPosition(1,1,0)`