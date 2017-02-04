== General
* Replace any public structs with protos. This allows for the possibility of non-go interaction.
* Test how well the mouse and keyboard handlers deal with unusual events. 
Try unplugging mouse/keyboard. Using multiple mice. This isn't something we should need to worry about, but better to be safe.
* Gestures / touchpad support. Probably requires modifying goxjs/glfw
* Fullscreen toggle
* Initial loading screen. Particularly for the webgl version. It takes a while to load.
* Trailing camera. Based on player position, but has its own max speed and delay, so stays behind player a bit. 
Something like: https://docs.unity3d.com/ScriptReference/Vector3.SmoothDamp.html
* DirectionalCamera that follows player orientation (up on the screen is always the direction the player faces).
* Draw with specific layers. Right now everything is based on the order of draw calls. This could be pushed off onto 
users of gome, but it is likely to be a very common requirement, particularly for 2D games. Better to support it.
* http://www.gopherjs.org/ #Performance Tips
  * Consider switching everything to float64 as it's more efficient with gopherjs, only if web performance is an issue.
* Look into golang benchmarks 
  * https://golang.org/pkg/testing/
  * https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
* Is there a way to get physical screen dimensions? a 1080p phone should have a different display (larger font for 
example) compared to a desktop monitor.
* The TargetCamera's Orthographic and Perspective projections don't deal with zoom in the same fashion. This means
their documentation about being able to zoom in/out to certain size percentages is not accurate.
* using gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{}) to bind a "null" buffer at the end of each draw call would be a
safe thing to do to prevent using the wrong buffer at some point - but BindBuffer calls are expensive.
* gl.UseProgram() is called way too often. Keep track of current shader in my shader package so only need to call 
 UseProgram() when it's necessary. Similar to using gl.BindBuffer(gl.Buffer{}), it would be safer to use 
 gl.UseProgram(gl.Program{}), but this adds even more expense.
* Can we precompile shaders?
* Add handling, especially on the server side, for signals, like SIGTERM: https://gobyexample.com/signals
* Current code is limited to max uint16 indices per mesh. This can be increased if we use the OES_element_index_uint
 extension for WebGL, and change the gl.DrawElements call to use gl.UNSIGNED_INT
* The way axes are drawn isn't smart. They are simply two points which make a line. This seems reasonable, until
you take into account camera drawing distances. You want the axis line to be "infinite" in length, but if you put
the endpoints too far away, then they are culled and the entire axis isn't drawn.

== Graphical
* Add motion blur https://github.com/goxjs/example/tree/master/motionblur
* Resizing screen shouldn't cause everything to be black.
* Anisotropic filtering
* When moving, little objects like stars appear smaller/darker and brighten when not moving. Unsure what is causing this.
* Add support for finding and setting the center of a model. Right now the input mesh can be normalized so the scale is 
right, but if the center is in the bottom left corner of a cube, the rotation is totally off.

== Thread Safety
* Mouse/keyboard handler reads and writes.
* FPS tracker needs mutex protection for its reads and writes.

== Web
* How to cache so the whole client doesn't need to be re-downloaded each time. Keep in mind, a new version of the client
will require at least a partial update.
  * This gzip utility may be handy: https://github.com/NYTimes/gziphandler
* Ideally the client will change very infrequently, although cached assets may need to be updated. 

== Goxjs
* support WebGL Extensions, like anisotropic filtering.
  * In webgl, to load an extension: http://blog.tojicode.com/2012/03/anisotropic-filtering-in-webgl.html
    ```
    var ext = gl.getExtension("MOZ_EXT_texture_filter_anisotropic");
    gl.texParameterf(gl.TEXTURE_2D, ext.TEXTURE_MAX_ANISOTROPY_EXT, 4);
    ```
    while in OpenGL use glGetString to get a list of all available extensions, and then enable it:
    https://www.khronos.org/opengles/sdk/docs/man/xhtml/glGetString.xml
    ```
 	extensions := gl.GetString(gl.EXTENSIONS)
 	if strings.Contains(extensions, "GL_EXT_texture_filter_anisotropic") {
 		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, 4)
 	}
    ```

== Existing Bugs
* Holding a key, then click and hold on the title bar, and release the key. It becomes stuck in the pressed state
since the key-release wasn't caught.
* Scrolling is erratic in the web build. It jumps around, so it must be detecting something. I expect the values
are just different from the 1.0 per tick on the desktop. 

