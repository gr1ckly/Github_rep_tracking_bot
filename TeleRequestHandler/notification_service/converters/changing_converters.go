package converters

import (
	"Common"
	"fmt"
)

func ConvertChanging(pattern string, dto Common.ChangingDTO) string {
	return fmt.Sprintf(pattern, dto.Link, dto.Event, dto.Author, dto.Title, dto.UpdatedAt)
}
