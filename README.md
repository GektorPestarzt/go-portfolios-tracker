# НИЯУ МИФИ. Лабораторная работа №1. Мясниов Руслан, Б21-525. 2024

## Описание системы

Сервис для инвесторов. Пользователь создает портфель и прикрепляет к нему токены своих кошельков. Он может отправить запрос на обновленние данных, перейдя в
статус Process. Запрос попадает в очередь. После его обработки (получения ответа с серверов брокеров), в зависимости от результата статуи меняется на Success или Error.
Для получения данных из базы данных предусмотрен GET запрос.

Ожидаемая нагрузка на систему: 100 - 500 RPS

### Используемые технологии

- Go
- SQLite
- RabbitMQ
- Python

В качестве брокера был выбран RabbitMQ из-за гибкости в обработки сообщений, простоты, надёжности.

## Запуск и тестирование системы

1. Запуск контейнера:
   ```bash
   cd deploy
   docker-compose up
   ```
2. Запуск сервера (Да, сервер запускается не в виртуальной среде из-за "Since Docker is a US company, we must comply with US export control regulations." А на танцы с бубном уже времени не было):
   ```bash
   go run cmd/main/app.go
   go run cmd/consumer/consumer.go
   ```
3. Запук тестирования:
   ```bash
   locust -H http://127.0.0.1:1234
   ```

## Нагрузочное тестирование

![chart](assets/chart.png)

Из графика видно, что по мере увеличения числа пользователей до 100, нагрузка на систему постепенно возрастает.
Когда количество пользователей достигает этой отметки, нагрузка стабилизируется на уровне примерно 200 запросов в секунду,
при этом количество ошибок в работе сервиса остается минимальным. Данные показывают, что сервис способен эффективно справляться с такой нагрузкой в течение пяти минут.

## Заключение

В ходе проведенной лабораторной работы была разработана система для инвесторов.
Для обеспечения взаимодействия между компонентами системы использовался брокер сообщений RabbitMQ.
Этот подход повысил устойчивость и пропускную способность системы, что было подтверждено тестами на высокие нагрузки,
во время которых система выдерживала до 200 запросов в секунду.
