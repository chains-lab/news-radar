package models

type RequestQuery struct {
	Filters       map[string]QueryFilter // Ключ-значение для фильтрации (например, {"username": "john"})
	SortField     string                 // Поле сортировки (например, "created_at")
	SortAscending bool                   // Направление сортировки: true - по возрастанию, false - по убыванию
	Limit         int64                  // Количество записей для выборки
	Offset        int64                  // Смещение для пагинации
}

type QueryFilter struct {
	Type   string      // Тип фильтра: "strict", "soft", "num", "date"
	Method string      // Метод сравнения: для strict — "eq" (или пусто), для num/date — "gt", "lt", "gte", "lte", для soft — "regex"
	Value  interface{} // Значение для фильтрации
}
