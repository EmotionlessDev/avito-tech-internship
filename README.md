# Тестовое задание на стажировку в Авито
[![Build](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/build.yml)
[![code-quality](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/code-quality.yml/badge.svg?branch=main)](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/code-quality.yml)
[![Go Tests](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/tests.yml/badge.svg)](https://github.com/EmotionlessDev/avito-tech-internship/actions/workflows/tests.yml)

## 📦 Требования

Перед запуском убедитесь, что установлены:

- **Docker** и **Docker Compose**
- **Make**
- **golang-migrate** (CLI)

Установка migrate:
```bash
brew install golang-migrate
# или
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
```

## 🔧 Переменные окружения

Создайте `.env` на основе примера:

```bash
cp .env.example .env
```

Заполните необходимые значения.

## ▶️ Как запустить проект

### 1. Собрать контейнеры

```bash
make build
```

### 2. Запустить проект

```bash
make up
```

### 3. Применить миграции

```bash
make migrate-up
```

### 4. Просмотреть логи

```bash
make logs
```

### 5. Остановить контейнеры

```bash
make down
```

## 🧹 Линтер

Запуск статического анализа:

```bash
make lint
```

## 🔄 Миграции

Создать новую миграцию:

```bash
migrate create -ext sql -dir migrations -seq <name>
```

Откат миграций:

```bash
make migrate-down
```
