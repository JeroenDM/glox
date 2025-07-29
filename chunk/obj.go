package chunk

type ObjKind byte

const (
	OBJ_STRING ObjKind = iota
)

type Obj struct {
	kind ObjKind
}

type ObjString struct {
	obj    *Obj
	length uint
	s      []byte
}
