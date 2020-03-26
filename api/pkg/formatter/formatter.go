package formatter

import (
	"fmt"
	"strings"
)

// FormatSQL .
func FormatSQL(condition bool, fieldLen int, name string, concat string, builder *strings.Builder, last bool, doIfHappens func()) {
	defer func() {
		if last {
			builder.WriteByte(')')
		}
	}()
	if !condition {
		return
	}
	if fieldLen == 0 {
		builder.WriteByte('(')
	} else {
		builder.WriteString(fmt.Sprintf(" %s ", concat))
	}
	builder.WriteString(fmt.Sprintf("%s=$%d", name, fieldLen+1))
	doIfHappens()
}
