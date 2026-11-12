package converter

type Converter interface {
	Convert([]byte) ([]byte, error)
}
