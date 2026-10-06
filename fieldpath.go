package gofactory

import (
	"go/types"
	"strings"
)

// fieldPath is one path from a type down to a field reached only through
// by-value struct fields and embedding, ending at a protected type.
// An empty path means the type itself is protected.
type fieldPath struct {
	path  []string
	named *types.Named
}

func (p fieldPath) String() string {
	return strings.Join(p.path, ".")
}

// The result includes t itself, with an empty path, when t is protected.
func (d *detector) fieldPaths(t types.Type) []fieldPath {
	return d.fieldPathsVisited(t, map[types.Type]bool{})
}

func (d *detector) fieldPathsVisited(
	typ types.Type, visited map[types.Type]bool,
) []fieldPath {
	if cached, ok := d.fieldPathCache[typ]; ok {
		return cached
	}

	if visited[typ] {
		return nil
	}

	visited[typ] = true
	defer delete(visited, typ)

	var paths []fieldPath

	if named, ok := protectedNamed(typ); ok {
		paths = append(paths, fieldPath{named: named})
	}

	if strukt, ok := typ.Underlying().(*types.Struct); ok {
		for field := range strukt.Fields() {
			for _, sub := range d.fieldPathsVisited(field.Type(), visited) {
				paths = append(paths, fieldPath{
					path:  append([]string{field.Name()}, sub.path...),
					named: sub.named,
				})
			}
		}
	}

	d.fieldPathCache[typ] = paths

	return paths
}
