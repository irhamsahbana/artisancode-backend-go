package pkg

import (
	"strconv"
	"strings"
)

func ToPgVectorString(vec []float64) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range vec {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	}
	sb.WriteByte(']')
	return sb.String()
}
