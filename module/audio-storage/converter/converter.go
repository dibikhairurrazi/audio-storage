package converter

type Converter interface {
	Convert([]byte, string) ([]byte, error)
}
