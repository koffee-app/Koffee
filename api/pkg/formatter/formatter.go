package formatter

import (
	"fmt"
	"strings"
)

// FormatWhereQuery Helpful for building up a SELECT WHERE (...) query
func FormatWhereQuery(condition bool, fieldLen int, name string, concat string, builder *strings.Builder, last bool, doIfHappens func()) {
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

// Array returns a string like (key=$1 OR key=$2), nVALUES
func Array(n int, key string, array []string) (string, []interface{}) {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteByte('(')
	newArray := make([]interface{}, len(array))
	for idx := 0; idx < n; idx++ {
		if idx > 0 {
			stringBuilder.WriteString(" OR ")
		}
		stringBuilder.WriteString(fmt.Sprintf("%s=$%d", key, idx+1))
		newArray[idx] = array[idx]
	}
	stringBuilder.WriteByte(')')
	return stringBuilder.String(), newArray
}

// ArrayUint32 returns a string like (key=$1), nValues)
func ArrayUint32(n int, key string, array []uint32) (string, []interface{}) {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteByte('(')
	newArray := make([]interface{}, len(array))
	for idx := 0; idx < n; idx++ {
		if idx > 0 {
			stringBuilder.WriteString(" OR ")
		}
		stringBuilder.WriteString(fmt.Sprintf("%s=$%d", key, idx+1))
		newArray[idx] = array[idx]
	}
	stringBuilder.WriteByte(')')
	return stringBuilder.String(), newArray
}

func IfTrueAdd(str *strings.Builder, condition bool, key string, value interface{}, values []interface{}) []interface{} {
	if condition {
		len := len(values)
		if len > 0 {
			str.WriteByte(',')
		}
		str.WriteString(fmt.Sprintf("%s=$%d", key, len+1))
		values = append(values, value)
		return values
	}
	return values
}
