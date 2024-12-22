
# Веб-сервис для вычисления арифметических выражений

## Описание проекта
Данный проект представляет собой веб-сервис, который позволяет выполнять вычисления арифметических выражений. Пользователь отправляет запрос через HTTP POST, указывая выражение, и получает в ответе результат вычисления либо сообщение об ошибке.

---

## API

### URL

**POST** `/api/v1/calculate`

### Тело запроса
```json
{
    "expression": "ваше выражение"
}
```
Где `expression` — строка, содержащая арифметическое выражение.

### Ответы сервера

1. **Успех (200):**
    ```json
    {
        "result": "результат вычисления"
    }
    ```

2. **Ошибка 422 (Неверное выражение):**
    ```json
    {
        "error": "Expression is not valid"
    }
    ```

3. **Ошибка 500 (Внутренняя ошибка сервера):**
    ```json
    {
        "error": "Internal server error"
    }
    ```

---

## Примеры использования

### Пример успешного запроса:
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
    -H "Content-Type: application/json" \
    -d '{"expression": "2+2*2"}'
```
**Ответ:**
```json
{
    "result": 6.000000
}
```

### Пример ошибки 422 (некорректное выражение):
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
    -H "Content-Type: application/json" \
    -d '{"expression": "2+2*2a"}'
```
**Ответ:**
```json
{
    "error": "Expression is not valid"
}
```

### Пример ошибки 500:
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
    -H "Content-Type: application/json" \
    -d '{"expression": "1/0"}'
```
**Ответ:**
```json
{
    "error": "Internal server error"
}
```

---

## Запуск проекта

1. Склонируйте репозиторий:
    ```bash
    git clone https://github.com/Portalshik/calc-LMS.git
    cd calc-LMS
    ```

2. Убедитесь, что у вас установлена версия Go 1.18 или выше.

3. Запустите сервер:
    - **Для Linux:**
        ```bash
        go run cmd/calc/main.go
        ```
    - **Для Windows:**
        ```bash
        go run cmd/calc/main.go
        ```

4. Сервер будет доступен по адресу: `http://localhost:8080`.

---

## Требования

- Go версии 1.18 или выше.
- Установленный HTTP-клиент, например `curl`.

---

## Структура проекта
```
calc-LMS/
├── cmd/
│   └── calc/
│       └── main.go
├── go.mod
├── internal/
│   ├── api/v1/
│   │   └── api.go
│   ├── calculator/
│   │   └── calculator.go
│   └── server/
│       └── server.go
└── README.md
```
