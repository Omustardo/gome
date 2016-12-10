Gome (GOlang gaME engine)

### Terms
 * Handler: singleton that deals with some sort of input. For example, mouse.Handler can be used to access mouse 
 related information like mouse.Handler.Position(). Note that all handlers require initialization via their 
 Initialize function.
 * NRGBA: non-premultiplied RGBA color. In a premultiplied RGBA color, the RGB values are always between 0 and the 
 Alpha value, since they've been scaled down. For example, if your NRGBA color is solid red and half transparent 
 ```[1,0,0,0.5]```  then the premultiplied RGBA version is ```[0.5,0,0,0.5]```. I believe users will find NRGBA more 
 natural to work with so that's the default.