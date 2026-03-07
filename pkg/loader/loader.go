// Package loader provides functions for loading 3D model files.
package loader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// faceKey uniquely identifies a vertex by its OBJ attribute indices.
type faceKey struct {
	posIdx, uvIdx, normIdx int
}

// LoadOBJ parses a Wavefront OBJ file and returns the contained meshes.
// Supports positions (v), texture coordinates (vt), normals (vn), and
// triangulated face definitions (f). Only the first object/group is returned
// as a single Mesh.
func LoadOBJ(path string) ([]*geometry.Mesh, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader.LoadOBJ: %w", err)
	}
	defer f.Close()

	var positions []math.Vec3
	var uvs [][2]float64
	var normals []math.Vec3

	mesh := geometry.NewMesh()
	vertexCache := make(map[faceKey]int)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "v":
			if len(fields) < 4 {
				continue
			}
			p, parseErr := parseVec3(fields[1], fields[2], fields[3])
			if parseErr != nil {
				return nil, fmt.Errorf("loader.LoadOBJ: bad v line %q: %w", line, parseErr)
			}
			positions = append(positions, p)

		case "vt":
			if len(fields) < 3 {
				continue
			}
			u, parseErr := strconv.ParseFloat(fields[1], 64)
			if parseErr != nil {
				return nil, fmt.Errorf("loader.LoadOBJ: bad vt line %q: %w", line, parseErr)
			}
			v, parseErr := strconv.ParseFloat(fields[2], 64)
			if parseErr != nil {
				return nil, fmt.Errorf("loader.LoadOBJ: bad vt line %q: %w", line, parseErr)
			}
			uvs = append(uvs, [2]float64{u, v})

		case "vn":
			if len(fields) < 4 {
				continue
			}
			n, parseErr := parseVec3(fields[1], fields[2], fields[3])
			if parseErr != nil {
				return nil, fmt.Errorf("loader.LoadOBJ: bad vn line %q: %w", line, parseErr)
			}
			normals = append(normals, n)

		case "f":
			// Triangulate fan for faces with >3 verts (n-2 triangles from vertex 0).
			if len(fields) < 4 {
				continue
			}
			fan := make([]int, 0, len(fields)-1)
			for _, token := range fields[1:] {
				idx, parseErr := resolveVertex(token, positions, uvs, normals, mesh, vertexCache)
				if parseErr != nil {
					return nil, fmt.Errorf("loader.LoadOBJ: bad f token %q: %w", token, parseErr)
				}
				fan = append(fan, idx)
			}
			for i := 1; i < len(fan)-1; i++ {
				mesh.AddTriangle(fan[0], fan[i], fan[i+1])
			}
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("loader.LoadOBJ: scanner error: %w", scanErr)
	}
	return []*geometry.Mesh{mesh}, nil
}

// resolveVertex maps an OBJ face token (e.g. "3/2/1") to a geometry.Vertex index,
// creating a new vertex in the mesh if this combination has not been seen.
func resolveVertex(
	token string,
	positions []math.Vec3,
	uvs [][2]float64,
	normals []math.Vec3,
	mesh *geometry.Mesh,
	cache map[faceKey]int,
) (int, error) {
	pi, ti, ni, err := parseFaceToken(token)
	if err != nil {
		return 0, err
	}

	// Convert 1-based OBJ indices to 0-based, handling optional components.
	posIdx := -1
	if pi != 0 {
		posIdx = pi - 1
	}
	uvIdx := -1
	if ti != 0 {
		uvIdx = ti - 1
	}
	normIdx := -1
	if ni != 0 {
		normIdx = ni - 1
	}

	key := faceKey{posIdx, uvIdx, normIdx}
	if existing, ok := cache[key]; ok {
		return existing, nil
	}

	v := geometry.Vertex{}
	if posIdx >= 0 && posIdx < len(positions) {
		v.Position = positions[posIdx]
	}
	if uvIdx >= 0 && uvIdx < len(uvs) {
		v.UV = uvs[uvIdx]
	}
	if normIdx >= 0 && normIdx < len(normals) {
		v.Normal = normals[normIdx]
	}

	idx := mesh.AddVertex(v)
	cache[key] = idx
	return idx, nil
}

// parseFaceToken parses an OBJ face token like "1", "1/2", "1//3", or "1/2/3"
// and returns the position, UV, and normal indices (1-based; 0 = absent).
func parseFaceToken(token string) (posIdx, uvIdx, normIdx int, err error) {
	parts := strings.Split(token, "/")
	posIdx, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid position index %q: %w", parts[0], err)
	}
	if len(parts) > 1 && parts[1] != "" {
		uvIdx, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid UV index %q: %w", parts[1], err)
		}
	}
	if len(parts) > 2 && parts[2] != "" {
		normIdx, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid normal index %q: %w", parts[2], err)
		}
	}
	return posIdx, uvIdx, normIdx, nil
}

func parseVec3(xs, ys, zs string) (math.Vec3, error) {
	x, err := strconv.ParseFloat(xs, 64)
	if err != nil {
		return math.Vec3{}, err
	}
	y, err := strconv.ParseFloat(ys, 64)
	if err != nil {
		return math.Vec3{}, err
	}
	z, err := strconv.ParseFloat(zs, 64)
	if err != nil {
		return math.Vec3{}, err
	}
	return math.Vec3{X: x, Y: y, Z: z}, nil
}
