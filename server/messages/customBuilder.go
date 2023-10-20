package messages

import (
	"fmt"
	"strings"
)

type CustomBuilder struct {
	strings.Builder
}

func (b *CustomBuilder) WriteString(s string) {
	b.Builder.WriteString(s + "\n")
}

func (b *CustomBuilder) WriteFString(str string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(str, args...))
}
