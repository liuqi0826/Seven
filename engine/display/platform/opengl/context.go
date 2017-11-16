package opengl

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/liuqi0826/seven/engine/display/platform"
	"github.com/liuqi0826/seven/engine/utils"
	"github.com/liuqi0826/seven/events"
	"github.com/vulkan-go/glfw/v3.3/glfw"
)

const (
	ALWAYS        = "always"
	EQUAL         = "equal"
	GREATER       = "greater"
	GREATER_EQUAL = "greaterEqual"
	LESS          = "less"
	LESS_EQUAL    = "lessEqual"
	NEVER         = "never"
	NOT_EQUAL     = "notEqual"
)

const (
	BYTES_4 = "bytes4"
	FLOAT_1 = "float1"
	FLOAT_2 = "float2"
	FLOAT_3 = "float3"
	FLOAT_4 = "float4"
)

type Context struct {
	events.EventDispatcher

	config               *utils.Config
	window               *glfw.Window
	currentShaderProgram *Program3D

	ready bool
	debug bool

	clearRed   float32
	clearGreen float32
	clearBlue  float32
	clearAlpha float32

	depthMask       bool
	passCompareMode string
}

func (this *Context) Setup(config *utils.Config) error {
	var err error

	this.EventDispatcher.EventDispatcher(this)

	this.config = config
	this.debug = config.Debug

	glfw.WindowHint(glfw.Resizable, glfw.False)
	//glfw.WindowHint(glfw.ContextVersionMajor, 4)
	//glfw.WindowHint(glfw.ContextVersionMinor, 5)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	this.window, err = glfw.CreateWindow(this.config.WindowWidth, this.config.WindowHeight, this.config.WindowTitle, nil, nil)
	if err != nil {
		return err
	}
	this.window.MakeContextCurrent()

	this.window.SetKeyCallback(this.keyboardCallback)
	this.window.SetMouseButtonCallback(this.mouseButtonCallback)
	this.window.SetCursorPosCallback(this.cursorPositionCallback)
	this.window.SetSizeCallback(this.resizeCallback)

	err = gl.Init()
	if err != nil {
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.CullFace(gl.BACK)

	this.ready = true

	return err
}
func (this *Context) GetWindow() *glfw.Window {
	return this.window
}
func (this *Context) Clear(color bool, depth bool, stencil bool) {
	var mask uint32
	if color {
		mask = mask | gl.COLOR_BUFFER_BIT
	}
	if depth {
		mask = mask | gl.DEPTH_BUFFER_BIT
	}
	if stencil {
		mask = mask | gl.STENCIL_BUFFER_BIT
	}
	if this.ready {
		gl.Clear(mask)
	}
}
func (this *Context) ConfigureBackBuffer() {
}
func (this *Context) CreateCubeTexture() {
}
func (this *Context) CreateProgram() platform.IProgram3D {
	return new(Program3D)
}
func (this *Context) CreateTexture() {
}
func (this *Context) CreateIndexBuffer() platform.IIndexBuffer {
	indexBuffer := new(IndexBuffer)
	gl.GenBuffers(1, &indexBuffer.Index)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer.Index)
	return indexBuffer
}
func (this *Context) CreateVertexBuffer() platform.IVertexBuffer {
	vertexBuffer := new(VertexBuffer)
	gl.GenBuffers(1, &vertexBuffer.Index)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.Index)
	return vertexBuffer
}
func (this *Context) Dispose() {
}
func (this *Context) Present() {
	this.window.SwapBuffers()
}
func (this *Context) DrawTriangles(indexBuffer platform.IIndexBuffer, firstIndex int32, numTriangles int32) {
	if indexBuf, ok := indexBuffer.(*IndexBuffer); ok {
		gl.BindVertexArray(indexBuf.Index)
		gl.DrawElements(gl.TRIANGLES, numTriangles, gl.UNSIGNED_SHORT, nil)
		gl.BindVertexArray(0)
	}
}
func (this *Context) SetBlendFactors() {
}
func (this *Context) SetClearColor(red float32, green float32, blue float32, alpha float32) {
	this.clearRed = red
	this.clearGreen = green
	this.clearBlue = blue
	this.clearAlpha = alpha
	if this.ready {
		gl.ClearColor(this.clearRed, this.clearGreen, this.clearBlue, this.clearAlpha)
	} else {
		fmt.Println("GL is not ready.")
	}
}
func (this *Context) SetColorMask() {
}
func (this *Context) SetCulling() {
}
func (this *Context) SetDepthTest(depthMask bool, passCompareMode string) {
	this.depthMask = depthMask
	this.passCompareMode = passCompareMode
	if this.ready {
		if this.depthMask {
			gl.Enable(gl.DEPTH_TEST)
		} else {

		}
		switch this.passCompareMode {
		case ALWAYS:
			gl.DepthFunc(gl.ALWAYS)
		case EQUAL:
			gl.DepthFunc(gl.EQUAL)
		case GREATER:
			gl.DepthFunc(gl.GREATER)
		case GREATER_EQUAL:
		case LESS:
			gl.DepthFunc(gl.LESS)
		case LESS_EQUAL:
		case NEVER:
			gl.DepthFunc(gl.NEVER)
		case NOT_EQUAL:
		}
	}
}
func (this *Context) SetProgram(program platform.IProgram3D) {
	if program3D, ok := program.(*Program3D); ok {
		this.currentShaderProgram = program3D
		gl.UseProgram(program3D.Index)
	}
}
func (this *Context) SetProgramConstantsFromByteArray() {
}
func (this *Context) SetProgramConstantsFromMatrix() {
}
func (this *Context) SetProgramConstantsFromVector() {
}
func (this *Context) SetRenderToBackBuffer() {
}
func (this *Context) SetRenderToTexture() {
}
func (this *Context) SetScissorRectangle() {
}
func (this *Context) SetStencilActions() {
}
func (this *Context) SetStencilReferenceValue() {
}
func (this *Context) SetTextureAt() {
}
func (this *Context) SetVertexBufferAt(value string, stride int32, bufferOffset int, format string) {
	if this.currentShaderProgram != nil {
		var size int32
		var xtype uint32

		switch format {
		case FLOAT_1:
			size = 1
			xtype = gl.FLOAT
		case FLOAT_2:
			size = 2
			xtype = gl.FLOAT
		case FLOAT_3:
			size = 3
			xtype = gl.FLOAT
		case FLOAT_4:
			size = 4
			xtype = gl.FLOAT
		case BYTES_4:
			size = 4
			xtype = gl.BYTE
		}

		attrib := uint32(gl.GetAttribLocation(this.currentShaderProgram.Index, gl.Str(value+"\x00")))
		gl.EnableVertexAttribArray(attrib)
		gl.VertexAttribPointer(attrib, size, xtype, false, stride, gl.PtrOffset(bufferOffset))
	}
}
func (this *Context) ShouldClose() bool {
	return this.window.ShouldClose()
}

func (this *Context) keyboardCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		return
	}
	fmt.Println(key, scancode, action, mods)
}
func (this *Context) mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	fmt.Println(button, action, mod)
}
func (this *Context) cursorPositionCallback(window *glfw.Window, x float64, y float64) {
	fmt.Println(x, y)
}
func (this *Context) frameBufferSizeCallback() {
}
func (this *Context) scrollCallback() {
}
func (this *Context) makeContextCurrentCallback() {
}
func (this *Context) resizeCallback(window *glfw.Window, width int, height int) {
	fmt.Println(width, height)
}
