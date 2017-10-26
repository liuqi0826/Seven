package display

import (
	//"fmt"
	"math"

	"github.com/liuqi0826/seven/engine/display/base"
	"github.com/liuqi0826/seven/engine/display/core"
	"github.com/liuqi0826/seven/engine/display/pick"
	"github.com/liuqi0826/seven/engine/static"
	"github.com/liuqi0826/seven/geom"
)

const PROJECTION_TYPE_ORTHO = "projectionTypeOrtho"
const PROJECTION_TYPE_PERSPECTIVE = "projectionTypePerspective"

const COORDINATE_SYSTEM_LEFT_HAND = "coordinateSystemLeftHand"
const COORDINATE_SYSTEM_RIGHT_HAND = "coordinateSystemRightHand"

type ProjectionConfig struct {
	ProjectType      string
	Width            float32
	Height           float32
	NearClipping     float32
	FarClipping      float32
	Field            float32
	CoordinateSystem string
}

type Camera struct {
	base.Object

	Controller  core.IController
	MousePicker pick.Picker
	DisplayList []core.IDisplayObject

	config *ProjectionConfig

	host *Scene

	upAxis       geom.Vector4
	projection   geom.Matrix4x4
	invTransform geom.Matrix4x4
}

func (this *Camera) Camera(host *Scene, config *ProjectionConfig) {
	this.Object.Object()

	this.host = host
	this.config = config
	if this.config == nil {
		this.config = new(ProjectionConfig)
		this.config.ProjectType = PROJECTION_TYPE_PERSPECTIVE
		this.config.Width = 1280
		this.config.Height = 720
		this.config.NearClipping = 0.1
		this.config.FarClipping = 1000.0
		this.config.Field = 45.0
		this.config.CoordinateSystem = COORDINATE_SYSTEM_LEFT_HAND
	}

	this.MousePicker.Picker()
	this.createProjectionMatrix()
}
func (this *Camera) Update() {
	this.Object.Update()
	if this.Controller != nil {
		this.Controller.Update()
	}
	if this.host != nil {
		this.DisplayList = this.host.displayList
	}

	this.invTransform = this.Object.GetTransform().Clone()
	this.invTransform.Invert()
}
func (this *Camera) GetProjectionMatrix() *geom.Matrix4x4 {
	return &this.projection
}
func (this *Camera) GetTransformMatrix() *geom.Matrix4x4 {
	return &this.invTransform
}
func (this *Camera) LookAt(at *geom.Vector4, up *geom.Vector4) {
	if up == nil {
		this.upAxis = geom.VERCTOR4_Y_AXIS.Clone()
	} else {
		this.upAxis = up.Clone()
	}
	switch this.config.CoordinateSystem {
	case COORDINATE_SYSTEM_LEFT_HAND:
		zAxis := at.Clone()
		zAxis.Subtract(this.GetPosition())
		zAxis.Normalize()
		xAxis := this.upAxis.Clone()
		xAxis.CrossProduct(&zAxis)
		xAxis.Normalize()
		yAxis := zAxis.Clone()
		yAxis.CrossProduct(&xAxis)
		xm := -xAxis.DotProduct(this.GetPosition())
		ym := -yAxis.DotProduct(this.GetPosition())
		zm := -zAxis.DotProduct(this.GetPosition())
		raw := [16]float32{
			xAxis.X, yAxis.X, zAxis.X, 0.0,
			xAxis.Y, yAxis.Y, zAxis.Y, 0.0,
			xAxis.Z, yAxis.Z, zAxis.Z, 0.0,
			xm, ym, zm, 1.0}
		mtx := new(geom.Matrix4x4)
		mtx.Matrix4x4(&raw)
		this.GetTransform().Append(mtx)
	case COORDINATE_SYSTEM_RIGHT_HAND:
		zAxis := this.GetPosition().Clone()
		zAxis.Subtract(at)
		zAxis.Normalize()
		xAxis := this.upAxis.Clone()
		xAxis.CrossProduct(&zAxis)
		xAxis.Normalize()
		yAxis := zAxis.Clone()
		yAxis.CrossProduct(&xAxis)
		xm := -xAxis.DotProduct(this.GetPosition())
		ym := -yAxis.DotProduct(this.GetPosition())
		zm := -zAxis.DotProduct(this.GetPosition())
		raw := [16]float32{
			xAxis.X, yAxis.X, zAxis.X, 0.0,
			xAxis.Y, yAxis.Y, zAxis.Y, 0.0,
			xAxis.Z, yAxis.Z, zAxis.Z, 0.0,
			xm, ym, zm, 1.0}
		mtx := new(geom.Matrix4x4)
		mtx.Matrix4x4(&raw)
		this.GetTransform().Append(mtx)
	}
}
func (this *Camera) createProjectionMatrix() {
	var raw [16]float32
	switch this.config.ProjectType {
	case PROJECTION_TYPE_ORTHO:
		switch this.config.CoordinateSystem {
		case COORDINATE_SYSTEM_LEFT_HAND:
			raw = [16]float32{
				2.0 / float32(this.config.Height), 0.0, 0.0, 0.0,
				0.0, 2.0 / float32(this.config.Height), 0.0, 0.0,
				0.0, 0.0, 1.0 / (this.config.FarClipping - this.config.NearClipping), 0.0,
				0.0, 0.0, this.config.NearClipping / (this.config.NearClipping - this.config.FarClipping), 1.0}
		case COORDINATE_SYSTEM_RIGHT_HAND:
			raw = [16]float32{
				2.0 / float32(this.config.Width), 0.0, 0.0, 0.0,
				0.0, 2.0 / float32(this.config.Height), 0.0, 0.0,
				0.0, 0.0, 1.0 / (this.config.NearClipping - this.config.FarClipping), 0.0,
				0.0, 0.0, this.config.NearClipping / (this.config.NearClipping - this.config.FarClipping), 1.0}
		}
	case PROJECTION_TYPE_PERSPECTIVE:
		aspectRatio := float32(this.config.Width) / float32(this.config.Height)
		yScale := 1.0 / float32(math.Tan(float64(this.config.Field/2.0)))
		xScale := yScale / aspectRatio
		switch this.config.CoordinateSystem {
		case COORDINATE_SYSTEM_LEFT_HAND:
			raw = [16]float32{
				xScale, 0.0, 0.0, 0.0,
				0.0, yScale, 0.0, 0.0,
				0.0, 0.0, this.config.FarClipping / (this.config.FarClipping - this.config.NearClipping), 1.0,
				0.0, 0.0, (this.config.NearClipping * this.config.FarClipping) / (this.config.NearClipping - this.config.FarClipping), 0.0}
		case COORDINATE_SYSTEM_RIGHT_HAND:
			raw = [16]float32{
				xScale, 0.0, 0.0, 0.0,
				0.0, yScale, 0.0, 0.0,
				0.0, 0.0, this.config.FarClipping / (this.config.NearClipping - this.config.FarClipping), -1.0,
				0.0, 0.0, (this.config.NearClipping * this.config.FarClipping) / (this.config.NearClipping - this.config.FarClipping), 0.0}
		}
	}
	this.projection = geom.Matrix4x4{}
	this.projection.Matrix4x4(&raw)
	switch static.API {
	case static.GL:
		raw = [16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, 1.0, 0.0, 0.0,
			0.0, 0.0, 2.0, 0.0,
			0.0, 0.0, -1.0, 1.0}
		mtx := new(geom.Matrix4x4)
		mtx.Matrix4x4(&raw)
		this.projection.Append(mtx)
	case static.VULKAN:
	case static.D3D9:
	case static.D3D12:
	}
}
