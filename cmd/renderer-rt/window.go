//go:build !headless

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
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

// run initialises a GLFW window and enters the CPU-blit render loop.
// It returns when the window is closed or ESC is pressed.
func run(cfg Config) error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw init: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(cfg.Width, cfg.Height, "go_toy_renderer", nil, nil)
	if err != nil {
		return fmt.Errorf("create window: %w", err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	if err = gl.Init(); err != nil {
		return fmt.Errorf("opengl init: %w", err)
	}

	prog, err := buildShaderProgram(vertSrc, fragSrc)
	if err != nil {
		return err
	}
	defer gl.DeleteProgram(prog)

	// VAO required by OpenGL 4.1 core profile (no default VAO).
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	defer gl.DeleteVertexArrays(1, &vao)

	tex := createTexture(cfg.Width, cfg.Height)
	defer gl.DeleteTextures(1, &tex)

	fb := framebuffer.New(cfg.Width, cfg.Height)
	scene := render.NewScene()
	scene.AddMesh(geometry.NewCube())
	cam := camera.New(
		math.Vec3{X: 3, Y: 2, Z: 5},
		math.Vec3{},
		math.Vec3{X: 0, Y: 1},
		45.0,
		float64(cfg.Width)/float64(cfg.Height),
		0.1,
		100.0,
	)

	gl.UseProgram(prog)
	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("uTex\x00")), 0)

	for !window.ShouldClose() {
		frameStart := time.Now()

		fb.Clear(math.Vec3{}, 1.0)
		render.Render(scene, cam, fb, shader.VertexColor)

		pixels := framebufferToRGBA(fb)
		uploadTexture(tex, cfg.Width, cfg.Height, pixels)

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		window.SwapBuffers()
		glfw.PollEvents()

		if elapsed := time.Since(frameStart); elapsed < frameDur {
			time.Sleep(frameDur - elapsed)
		}
	}

	return nil
}

// createTexture allocates an RGBA OpenGL texture sized width×height.
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
func buildShaderProgram(vertSrc, fragSrc string) (uint32, error) {
	vert, err := compileShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("vertex shader: %w", err)
	}
	defer gl.DeleteShader(vert)

	frag, err := compileShader(fragSrc, gl.FRAGMENT_SHADER)
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
		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(prog, logLen, nil, gl.Str(log))
		gl.DeleteProgram(prog)
		return 0, fmt.Errorf("link program: %s", log)
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
		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(s, logLen, nil, gl.Str(log))
		gl.DeleteShader(s)
		return 0, fmt.Errorf("compile: %s", log)
	}
	return s, nil
}
