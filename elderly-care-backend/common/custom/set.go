package custom

type Set interface {
	Add(value interface{})
	Remove(value interface{})
	Contains(value interface{}) bool
	Clear()
	IsEmpty() bool
	Size() int
}

type HashSet map[interface{}]struct{}

func NewHashSet(values ...interface{}) Set {
	h := HashSet{}
	for _, value := range values {
		h.Add(value)
	}
	return h
}

func (h HashSet) Add(value interface{}) {
	//TODO implement me
	h[value] = struct{}{}
}

func (h HashSet) Remove(value interface{}) {
	//TODO implement me
	delete(h, value)
}

func (h HashSet) Clear() {
	//TODO implement me
	clear(h)
}

func (h HashSet) Contains(value interface{}) bool {
	//TODO implement me
	_, ok := h[value]
	return ok
}

func (h HashSet) IsEmpty() bool {
	//TODO implement me
	return len(h) == 0
}

func (h HashSet) Size() int {
	//TODO implement me
	return len(h)
}
