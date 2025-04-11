package bot

func GetCommandsDescription() map[string]string {
	return map[string]string{
		"help":  "Справка о доступных командах",
		"start": "Начало работы с ботом",
		"add":   "Начало отслеживания нового репозитория",
		"del":   "Прекращение отслеживания репозитория",
		"repos": "Вывод отслеживаемых репозиториев",
	}
}
