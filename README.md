# Music-library
Music library REST full API

## Описание
REST full API для управления библиотекой песен с возможностью:
- Добавления/удаления/редактирования треков
- Получения песен с пагинацией
- Получения текстов песен с пагинацией
- Интеграции с внешними API

## Стек
* Язык: Go (Echo Framework)
* База данных: PostgreSQL
* Документация: Swagger/OpenAPI 3.0
* Контейнеризация: docker-compose

# Запуск
1. Клонирование репозитория:
git clone https://github.com/Dmitrii30002/Music-library.git
cd Music_library

2. Настройка конфигурации:
В файле .env находятся данные кофнигурации, при необходимости поменяйте их.

3. Контейнеризация базы данных:
make docker-compose-up

4. Запуск проекта:
make run

# Эндпоинты:

|Метод	    |Путь	              |Описание                                    |
|:----------|:------------------|:-------------------------------------------|
|GET	      |/songs	            |Список всех песен (с фильтрами и пагинацией)|
|POST	      |/songs	            |Добавить новую песню                        |
|PUT	      |/songs/:id	        |Обновить данные песни                       |
|DELETE	    |/songs/:id	        |Удалить песню                               |
|GET	      |/songs/:id/lyrics	|Получить текст с пагинацией                 |

