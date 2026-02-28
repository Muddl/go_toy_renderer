//go:build !headless

package renderer

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	pkgmath "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

const (
	targetFPS = 60
	frameDur  = time.Second / targetFPS
)

// Minimal vertex shader: pass through a fullscreen clip-space quad via gl_VertexID.
var vertSrc = `#version 410 core
out vec2 vUV;
void main() {
    vec2 pos[4] = vec2[4](
        vec2(-1.0,  1.0),
        vec2(-1.0, -1.0),
        vec2( 1.0,  1.0),
        vec2( 1.0, -1.0)
    );
    vec2 uv[4] = vec2[4](
        vec2(0.0, 0.0),
        vec2(0.0, 1.0),
        vec2(1.0, 0.0),
        vec2(1.0, 1.0)
    );
    vUV = uv[gl_VertexID];
    gl_Position = vec4(pos[gl_VertexID], 0.0, 1.0);
}
` + "\x00"

// Fragment shader: sample the CPU framebuffer texture.
var fragSrc = `#version 410 core
in vec2 vUV;
out vec4 fragColor;
uniform sampler2D uTex;
void main() {
    fragColor = texture(uTex, vUV);
}
` + "\x00"

// CPUBackend renders each frame on the CPU and blits the result via OpenGL.
type CPUBackend struct {
	width, height int
	window        *glfw.Window
	prog          uint32
	vao           uint32
	tex           uint32
	fb            *framebuffer.Framebuffer
	cam           camera.Camera
}

// Init opens a GLFW window, initialises OpenGL, and prepares the CPU framebuffer.
func (c *CPUBackend) Init(width, height int) error {
	c.width = width
	c.height = height

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw init: %w", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(width, height, "go_toy_renderer", nil, nil)
	if err != nil {
		glfw.Terminate()
		return fmt.Errorf("create window: %w", err)
	}
	c.window = win
	win.MakeContextCurrent()

	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	if err = gl.Init(); err != nil {
		glfw.Terminate()
		return fmt.Errorf("opengl init: %w", err)
	}

	prog, err := buildShaderProgram(vertSrc, fragSrc)
	if err != nil {
		glfw.Terminate()
		return err
	}
	c.prog = prog

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	c.vao = vao

	c.tex = createTexture(width, height)
	c.fb = framebuffer.New(width, height)
	c.cam = camera.New(
		pkgmath.Vec3{X: 3, Y: 2, Z: 5},
		pkgmath.Vec3{},
		pkgmath.Vec3{X: 0, Y: 1},
		45.0,
		float64(width)/float64(height),
		0.1,
		100.0,
	)

	gl.UseProgram(c.prog)
	gl.Uniform1i(gl.GetUniformLocation(c.prog, gl.Str("uTex\x00")), 0)

	return nil
}

// RenderFrame renders one frame of scene and presents it in the window.
// Returns ErrWindowClosed when the window has been dismissed.
func (c *CPUBackend) RenderFrame(scene *render.Scene) error {
	if c.window.ShouldClose() {
		return ErrWindowClosed
	}

	frameStart := time.Now()

	render.Render(scene, c.cam, c.fb, shader.VertexColor)

	pixels := framebufferToRGBA(c.fb)
	uploadTexture(c.tex, c.width, c.height, pixels)

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindVertexArray(c.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	c.window.SwapBuffers()
	glfw.PollEvents()

	if elapsed := time.Since(frameStart); elapsed < frameDur {
		time.Sleep(frameDur - elapsed)
	}

	return nil
}

// Shutdown releases OpenGL resources and terminates GLFW.
func (c *CPUBackend) Shutdown() {
	if c.tex != 0 {
		gl.DeleteTextures(1, &c.tex)
	}
	if c.vao != 0 {
		gl.DeleteVertexArrays(1, &c.vao)
	}
	if c.prog != 0 {
		gl.DeleteProgram(c.prog)
	}
	glfw.Terminate()
}

// createTexture allocates an RGBA OpenGL texture sized width x height.
func createTexture(width, height int) uint32 {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8,
		int32(width), int32(height), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, nil)
	return tex
}

// uploadTexture writes pixels into an already-allocated texture via TexSubImage2D.
func uploadTexture(tex uint32, width, height int, pixels []byte) {
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0,
		int32(width), int32(height),
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(pixels))
}

// buildShaderProgram compiles vertex and fragment GLSL sources and links them.
func buildShaderProgram(vSrc, fSrc string) (uint32, error) {
	vert, err := compileShader(vSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("vertex shader: %w", err)
	}
	defer gl.DeleteShader(vert)

	frag, err := compileShader(fSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("fragment shader: %w", err)
	}
	defer gl.DeleteShader(frag)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLen)
		logBuf := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(prog, logLen, nil, gl.Str(logBuf))
		gl.DeleteProgram(prog)
		return 0, fmt.Errorf("link program: %s", logBuf)
	}
	return prog, nil
}

// compileShader compiles a GLSL shader of the given type.
func compileShader(src string, shaderType uint32) (uint32, error) {
	s := gl.CreateShader(shaderType)
	csrc, free := gl.Strs(src)
	gl.ShaderSource(s, 1, csrc, nil)
	free()
	gl.CompileShader(s)

	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLen)
		logBuf := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(s, logLen, nil, gl.Str(logBuf))
		gl.DeleteShader(s)
		return 0, fmt.Errorf("compile: %s", logBuf)
	}
	return s, nil
}
