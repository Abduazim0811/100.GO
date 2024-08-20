package embed

import _ "embed"

//go:embed templates
var embedded []byte

func GetTemplate() []byte {
	return embedded
}
