# Система скачивания файлов zip-архивами.

## Описание
Система позволяет создавать задачи на формирование ZIP-архивов из файлов, доступных по URL. Поддерживаются только PDF и JPEG файлы. Максимальное количество файлов в одном архиве - 3. Одновременно может обрабатываться не более 3 задач.

## Api endpoints

### 1. Создание задачи
`POST /tasks`
Возвращает: `{"id": "<task_id>"}`
Статусы:
- 201: Задача создана
- 503: Сервер перегружен

### 2. Добавление URL в задачу
`POST /tasks/{id}/urls`
Параметры:
`url`: URL файла
Статусы:
- 200: URL добавлен
- 400: Неверный запрос
- 403: Достигнут лимит файлов
- 404: Задача не найдена

### 3. Получение статуса задачи
`GET /tasks/{id}`
Возвращает:
  ```json
  {
    "id": "task_id",
    "status": "PENDING|PROCESSING|COMPLETED|FAILED",
    "urls": [links_to_files],
    "archive": "download_url",
    "created_at": "date_time"
  }
```

## Особенности реализации

1. Пакетная структура и MVC
2. State Pattern, Синхронизация, Фоновая обработка
3. Обработка HTTP-запросов
4. Yaml-конфиг и инииализация состояния
5. Обработка ошибок и статусов
6. Facade Pattern

## Запуск

Файл main.go в папке cmd запускает приложение.

``` bash
go get -u
// из папки проекта
go run cmd/main.go
```
### Сторониие пакеты
``` bash
go get gopkg.in/yaml.v3
go get github.com/gorilla/mux
```

## Тестирование

### Командная строка
1. Создание задачи:
  ``` bash
   curl -X POST http://localhost:8080/tasks
  ```
Возвращает id задачи, на которое в последующих командах заменяется {task_id}.

2. Добавление ссылки на файл:

Ссылка на файл вместо {paste_url}
   ``` bash
    curl -X POST -d "url={paste_url}" http://localhost:8080/tasks/{task_id}/urls
   ```
3. Проверка статуса:
   ``` bash
    curl http://localhost:8080/tasks/{task_id}
   ```

### Postman
1. POST http://localhost:8080/tasks
2. POST http://localhost:8080/tasks/{task_id}/urls
3. GET http://localhost:8080/tasks/{task_id}
