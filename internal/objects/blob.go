package objects

func NewBlob(data []byte) object {
	return New(ObjectTypeBlob, data)
}
