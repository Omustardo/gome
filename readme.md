Gome (GOlang gaME engine)

### Terms
 * NRGBA: non-premultiplied RGBA color. In a premultiplied RGBA color, the RGB values are always between 0 and the 
 Alpha value, since they've been scaled down. For example, if your NRGBA color is solid red and half transparent 
 ```[1,0,0,0.5]```  then the premultiplied RGBA version is ```[0.5,0,0,0.5]```. I believe users will find NRGBA more 
 natural to work with so that's the default. Note that you should access the color values directly, i.e. `color.R`, 
 rather than using the `RGBA()` function.
 * Mesh: a bunch of 3d points are read in and stored in buffers on the GPU, and references to these buffers are put
 into mesh objects.
 * Entity: information about an object in the world. Position, Rotation, Scale. 
 * Model: The combination of a Mesh and Entity. Uses the world information from the entity, and the graphics information
 from the mesh to render.
 
###
 * The scale of models in the game is up to the user to set via entity.Entity's Scale field. Note that leaving it empty
 will mean your model has no physical size and so won't be shown. If you want your model to not be rendered, it's 
 more effective to set model.Hidden = true
 * Be careful when using mgl32. Its functions tend to return new structs rather than modifying. For example, you might
 be tempted to modify an Entity's position using ```e.Position.Add(mgl32.Vec3{1,1,0})``` but this doesn't modify 
 ```e.Position```. It just returns a new vector. Instead, use ```e.Position = e.Position.Add(mgl32.Vec3{1,1,0})```
  or ```e.ModifyCenter(1,1,0)```