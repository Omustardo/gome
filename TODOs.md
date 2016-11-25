== General
* Test how well the mouse and keyboard handlers deal with unusual events. 
Try unplugging mouse/keyboard. Using multiple mice.
* Gestures / touchpad support. Probably requires modifying goxjs/glfw
* Fullscreen toggle
* All of the game logic needs to be based on delta time since it was last applied. Right now it's based on happening 
per-frame which isn't consistent, and definitely won't work for multiplayer.
* Initial loading screen. Particularly for the webgl version. It takes a while to load.
* Trailing camera. Based on player position, but has its own max speed and delay, so stays behind player a bit. 
Something like: https://docs.unity3d.com/ScriptReference/Vector3.SmoothDamp.html
* DirectionalCamera that follows player orientation (up on the screen is always the direction the player faces).
* Draw with specific layers. Right now everything is based on the order of draw calls.
* http://www.gopherjs.org/ #Performance Tips
  * Consider switching everything to float64 as it's more efficient with gopherjs if performance is an issue.

== Thread Safety
* Mouse/keyboard handler reads and writes.
* FPS tracker needs mutex protection for its reads and writes.

== Web
* How to cache so the whole client doesn't need to be re-downloaded each time. Keep in mind, a new version of the client
will require at least a partial update.

== Bugs
* Holding a key, then click and hold on the title bar, and release the key. It becomes stuck in the pressed state
since the key-release wasn't caught.
* Scrolling is erratic in the web build. It jumps around, so it must be detecting something. I expect the values
are just different from the 1.0 per tick on the desktop. 