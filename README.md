# logcheck — линтер для проверки лог-сообщений (golangci-lint plugin)

Проект реализует кастомный линтер (анализатор на базе `go/analysis`), совместимый с **golangci-lint v2 module plugins**.

Проверяемые правила (по ТЗ):
1. Сообщение начинается со строчной буквы
2. Сообщение только на английском (ASCII)
3. Сообщение не содержит спецсимволов/эмодзи (по умолчанию разрешены только буквы/цифры/пробелы)
4. Сообщение не содержит потенциально чувствительных данных (по ключевым словам)

Поддерживаемые логгеры:
- `log/slog`
- `go.uber.org/zap` (в т.ч. sugared варианты `Infof/Infow/...`)

---

## Требования

- Go **1.22+**
- golangci-lint **v2+**

> ⚠️ На Windows `-buildmode=plugin` не поддерживается, собирайте плагин в **WSL** или на Linux/macOS.

---

## Тесты

Запуск всех тестов:

```bash
go test ./...
```

---

## Сборка плагина

```bash
go build -buildmode=plugin -o logcheck.so ./plugin
```

---

## Запуск через golangci-lint

Проект содержит файл `.golangci-plugin.yml` (для golangci-lint v2 module plugins).

Пример запуска (из корня проекта):

```bash
golangci-lint run -c .golangci-plugin.yml ./...
```

---

## Конфигурация (бонус)

Линтер умеет читать настройки из файла **.logcheck.yml** в рабочей директории.

Пример:

```yaml
# список ключевых слов для чувствительных данных
sensitive_keywords:
  - password
  - token
  - api_key

# какие спецсимволы разрешить (по умолчанию пусто, т.е. запрещены все спецсимволы)
allowed_special_chars: ""

# отключить конкретные правила
disabled_rules:
  - english-only
```

Доступные имена правил:
- `lowercase`
- `english-only`
- `no-special-chars`
- `no-sensitive-data`
