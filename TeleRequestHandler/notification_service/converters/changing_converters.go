package converters

import (
	"Common"
	"fmt"
)

const MESSAGE_PATTERN = "Репозиторий: %v\n" +
	"Событие: %v\n" +
	"Автор: %v\n" +
	"Комментарий: %v\n" +
	"Обновлено: %v"

func ConvertChanging(dto Common.ChangingDTO) string {
	return fmt.Sprintf(MESSAGE_PATTERN, dto.Link, dto.Event, dto.Author, dto.Title, dto.UpdatedAt)
}
