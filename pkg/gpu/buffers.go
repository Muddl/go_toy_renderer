//go:build !headless

package gpu

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
)

// VertexBuffer creates a GPU vertex buffer from vertices and uploads data via
// Queue.WriteBuffer. Buffer usage is Vertex|CopyDst. Returns an error if the
// Device has not been initialized.
func VertexBuffer(d *Device, vertices []geometry.Vertex) (*wgpu.Buffer, error) {
	if d.device == nil {
		return nil, errors.New("gpu: VertexBuffer: device not initialized")
	}
	data := float32SliceToBytes(PackVertices(vertices))
	size := uint64(len(data))
	if size == 0 {
		return nil, errors.New("gpu: VertexBuffer: no vertex data to upload")
	}
	buf := d.device.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  size,
		Usage: gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
	})
	if buf == nil {
		return nil, errors.New("gpu: VertexBuffer: CreateBuffer returned nil")
	}
	d.queue.WriteBuffer(buf, 0, data)
	return buf, nil
}

// IndexBuffer creates a GPU index buffer from integer indices and uploads data
// via Queue.WriteBuffer. Indices are stored as uint32 (IndexFormatUint32).
// Buffer usage is Index|CopyDst. Returns an error if the Device has not been
// initialized.
func IndexBuffer(d *Device, indices []int) (*wgpu.Buffer, error) {
	if d.device == nil {
		return nil, errors.New("gpu: IndexBuffer: device not initialized")
	}
	data := uint32SliceToBytes(PackIndices(indices))
	size := uint64(len(data))
	if size == 0 {
		return nil, errors.New("gpu: IndexBuffer: no index data to upload")
	}
	buf := d.device.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  size,
		Usage: gputypes.BufferUsageIndex | gputypes.BufferUsageCopyDst,
	})
	if buf == nil {
		return nil, errors.New("gpu: IndexBuffer: CreateBuffer returned nil")
	}
	d.queue.WriteBuffer(buf, 0, data)
	return buf, nil
}

// float32SliceToBytes converts []float32 to a little-endian []byte representation.
func float32SliceToBytes(data []float32) []byte {
	out := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(out[i*4:], math.Float32bits(v))
	}
	return out
}

// uint32SliceToBytes converts []uint32 to a little-endian []byte representation.
func uint32SliceToBytes(data []uint32) []byte {
	out := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(out[i*4:], v)
	}
	return out
}

// meshBufferEntry holds cached GPU buffers for a single mesh.
type meshBufferEntry struct {
	vertexBuf     *wgpu.Buffer
	vertexBufSize uint64
	indexBuf      *wgpu.Buffer
	indexBufSize  uint64
	indexCount    uint32
}

// LoadGeometry uploads mesh vertex and index data to VRAM and caches the result.
// If the same mesh pointer is passed again, the cached buffers are reused.
// Multiple meshes can be cached simultaneously so buffers remain valid during
// a render pass with multiple draw calls.
// Returns an error if the Device has not been initialized via Init.
func (d *Device) LoadGeometry(mesh *geometry.Mesh) error {
	if !d.ready {
		return errors.New("gpu: LoadGeometry: Device not initialized — call Init first")
	}
	if mesh == nil {
		return errors.New("gpu: LoadGeometry: mesh is nil")
	}

	if d.meshCache == nil {
		d.meshCache = make(map[*geometry.Mesh]*meshBufferEntry)
	}

	// Check cache hit.
	if entry, ok := d.meshCache[mesh]; ok {
		d.vertexBuf = entry.vertexBuf
		d.vertexBufSize = entry.vertexBufSize
		d.indexBuf = entry.indexBuf
		d.indexBufSize = entry.indexBufSize
		d.indexCount = entry.indexCount
		return nil
	}

	// Cache miss: upload new buffers.
	vb, err := VertexBuffer(d, mesh.Vertices)
	if err != nil {
		return fmt.Errorf("gpu: LoadGeometry: %w", err)
	}
	ib, err := IndexBuffer(d, mesh.Indices)
	if err != nil {
		vb.Destroy()
		vb.Release()
		return fmt.Errorf("gpu: LoadGeometry: %w", err)
	}

	entry := &meshBufferEntry{
		vertexBuf:     vb,
		vertexBufSize: uint64(len(mesh.Vertices)) * vertexStride,
		indexBuf:      ib,
		indexBufSize:  uint64(len(mesh.Indices)) * 4, // uint32 = 4 bytes
		indexCount:    uint32(len(mesh.Indices)),
	}
	d.meshCache[mesh] = entry

	d.vertexBuf = entry.vertexBuf
	d.vertexBufSize = entry.vertexBufSize
	d.indexBuf = entry.indexBuf
	d.indexBufSize = entry.indexBufSize
	d.indexCount = entry.indexCount
	return nil
}
