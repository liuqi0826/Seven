package display

import (
	"fmt"

	"github.com/liuqi0826/seven/engine/display/core"
	"github.com/liuqi0826/seven/events"
	"github.com/liuqi0826/seven/geom"
)

type Scene struct {
	events.EventDispatcher

	camera *Camera

	displayList []core.IDisplayObject
}

func (this *Scene) Scene() {
	this.EventDispatcher.EventDispatcher()

	this.camera = new(Camera)
	this.camera.Camera(this, nil)
	this.camera.X = 0
	this.camera.Y = 0
	this.camera.Z = -2

	zero := new(geom.Vector4)
	zero.Vector4()
	//this.Camera.LookAt(zero, nil)

	this.displayList = make([]core.IDisplayObject, 0)
}
func (this *Scene) AddChild(displayObject core.IDisplayObject) {
	fmt.Println(displayObject)
	this.displayList = append(this.displayList, displayObject)

	displayObject.SetRoot(this)
	displayObject.SetParent(this)
	displayObject.SetCamera(this.camera)

	event := new(events.Event)
	event.Type = events.ADDED
	displayObject.DispatchEvent(event)
}
func (this *Scene) RemoveChild(displayObject core.IDisplayObject) core.IDisplayObject {
	for i, c := range this.displayList {
		if c == displayObject {
			this.displayList = append(this.displayList[:i], this.displayList[i+1:]...)
			c.SetRoot(nil)
			c.SetParent(nil)
			c.SetCamera(nil)
			event := new(events.Event)
			event.Type = events.REMOVE
			c.DispatchEvent(event)
			return c
		}
	}
	return nil
}
func (this *Scene) RemoveAllChildren() {
	for _, c := range this.displayList {
		this.RemoveChild(c)
	}
	this.displayList = make([]core.IDisplayObject, 0)
}
func (this *Scene) GetChildByName(name string) core.IDisplayObject {
	for _, c := range this.displayList {
		if c.GetName() == name {
			return c
		}
	}
	return nil
}
func (this *Scene) GetCamera() *Camera {
	return this.camera
}
func (this *Scene) SetCamera(camera *Camera) {
	this.camera = camera
}
func (this *Scene) GetChildrenNumber() int32 {
	return int32(len(this.displayList))
}
