package resource

import (
	"github.com/liuqi0826/seven/geom"
)

type Geometry struct {
	ID       string
	Shader   string
	Material string

	Bound geom.Box

	VertexIndex []uint16
	ColorIndex  []uint16
	Vertex      [][]float32
}

func (this *Geometry) Geometry() {
	this.Bound = geom.Box{}
	this.VertexIndex = make([]uint16, 0)
	this.ColorIndex = make([]uint16, 0)
	this.Vertex = make([][]float32, 0)
}
