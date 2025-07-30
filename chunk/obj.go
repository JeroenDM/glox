package chunk

type ObjKind byte

const (
	OBJ_STRING ObjKind = iota
)

type Obj struct {
	kind ObjKind
}

type ObjString struct {
	Obj
	Length int
	Bytes  []byte
}

func CopyString(s []byte) ObjString {
	// TODO: or do we want null-terminated []byte arrays?
	dest := make([]byte, len(s))
	copy(dest, s)
	return ObjString{
		Obj:    Obj{kind: OBJ_STRING},
		Length: len(dest),
		Bytes:  dest,
	}
}

func TakeString(s []byte) ObjString {
	// TODO: or do we want null-terminated []byte arrays?
	return ObjString{
		Obj:    Obj{kind: OBJ_STRING},
		Length: len(s),
		Bytes:  s,
	}
}
