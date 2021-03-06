package immutable_map

type node struct {
	b     byte
	nodes Nodes
	value interface{}
}

func (a *node) insert(path []byte, value interface{}) *node {
	// override value
	if len(path) == 1 {
		return &node{
			b:     a.b,
			nodes: a.nodes,
			value: value,
		}
	}
	return &node{
		b:     a.b,
		nodes: a.nodes.insert(path[1:], value),
		value: a.value,
	}
}

func (a *node) contains(path []byte) bool {
	if len(path) == 0 {
		return true
	}
	return a.nodes.contains(path)
}

func (a *node) get(path []byte) (interface{}, bool) {
	return a.nodes.get(path)
}

func newNode(path []byte, value interface{}) *node {
	if len(path) == 1 {
		return &node{
			b:     path[0],
			nodes: Nodes{},
			value: value,
		}
	}
	return &node{
		b:     path[0],
		nodes: Nodes{}.insert(path[1:], value),
	}
}

type Nodes []*node

func (a Nodes) insert(path []byte, value interface{}) Nodes {
	if len(path) == 0 {
		return a
	}
	exists, index := findPosForInsert(a, path[0])
	clone := dup(a)
	if exists {
		clone[index] = clone[index].insert(path, value)
		return clone
	}
	clone = append(clone[:index], append(Nodes{newNode(path, value)}, clone[index:]...)...)
	return clone
}

func (a Nodes) contains(path []byte) bool {
	exists, index := contains(a, path[0])
	if !exists {
		return false
	}
	return a[index].contains(path[1:])
}

func (a Nodes) get(path []byte) (interface{}, bool) {
	exists, index := contains(a, path[0])
	if !exists {
		return nil, false
	}
	if len(path) == 1 {
		value := a[index].value
		return value, value != nil
	}
	return a[index].get(path[1:])
}

func dup(nodes []*node) []*node {
	out := make([]*node, len(nodes))
	for i, v := range nodes {
		out[i] = &*v
	}
	return out
}

func findPosForInsert(nodes []*node, b byte) (exists bool, pos int) {
	for i, v := range nodes {
		if b < v.b {
			return false, i
		}
		if b == v.b {
			return true, i
		}
		if b > v.b {
			continue
		}
	}
	return false, len(nodes)
}

func contains(nodes []*node, b byte) (exists bool, pos int) {
	for i, v := range nodes {
		if b == v.b {
			return true, i
		}
	}
	return false, len(nodes)
}

type Map struct {
	nodes Nodes
	count int
}

// Creates new Map.
func New() *Map {
	return &Map{}
}

func (a Map) Contains(path []byte) bool {
	if len(path) == 0 {
		return false
	}
	return a.nodes.contains(path)
}

// Insert value, uniquely identified by path bytes.
func (a *Map) Insert(path []byte, value interface{}) *Map {
	return &Map{
		nodes: a.nodes.insert(path, value),
	}
}

// Returns value identified by path bytes.
func (a Map) Get(path []byte) (interface{}, bool) {
	return a.nodes.get(path)
}

// Same as Get, but returns only one value.
func (a Map) Get1(path []byte) interface{} {
	rs, _ := a.Get(path)
	return rs
}

func (a Map) Count() int {
	return len(a.ToStringMap())
}

func (a Map) ToStringMap() map[string]interface{} {
	m := make(map[string]interface{})
	mapify(m, a.nodes, nil)
	return m
}

func mapify(m map[string]interface{}, n Nodes, path []byte) {
	for _, v := range n {
		if v.value != nil {
			m[string(append(path, v.b))] = v.value
		}
		mapify(m, v.nodes, append(path, v.b))
	}
}

func (a Map) ToSlice() []KeyValue {
	var m []KeyValue
	slicify(&m, a.nodes, nil)
	return m
}

func slicify(m *[]KeyValue, n Nodes, path []byte) {
	for _, v := range n {
		if v.value != nil {
			*m = append(*m, KeyValue{
				Key:   append(path, v.b),
				Value: v.value,
			})
		}
		slicify(m, v.nodes, append(path, v.b))
	}
}

type KeyValue struct {
	Key   []byte
	Value interface{}
}
