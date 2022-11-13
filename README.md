# Тестовое задание на позицию стажёра-бэкендера
## Микросервис для работы с балансом пользователей

### Задача:
Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 

### Требования к сервису:

1. Сервис должен предоставлять **HTTP API** с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: **Golang**.
2. Фреймворки и библиотеки можно использовать любые.
3. Реляционная СУБД: **MySQL** или **PostgreSQL**.
4. Использование docker и docker-compose для поднятия и развертывания dev-среды.
4. Весь код должен быть выложен на Github с **README** файлом с инструкцией по запуску и примерами запросов/ответов (можно просто описать в Readme методы, можно через Postman, можно в Readme curl запросы скопировать, и так далее).
5. Если есть потребность в асинхронных сценариях, то использование любых систем очередей - допускается.
6. При возникновении вопросов по ТЗ оставляем принятие решения за кандидатом (в таком случае Readme файле к проекту должен быть указан список вопросов с которыми кандидат столкнулся и каким образом он их решил).
7. Разработка интерфейса в браузере НЕ ТРЕБУЕТСЯ. Взаимодействие с API предполагается посредством запросов из кода другого сервиса. Для тестирования можно использовать любой удобный инструмент. Например: в терминале через curl или Postman.

### Будет плюсом:

[] Покрытие кода тестами.
[x] Swagger файл для вашего API. 
[x] Реализовать сценарий разрезервирования денег, если услугу применить не удалось.

### Дополнительные задания:

[x] **Доп. задание 1:**
Бухгалтерия раз в месяц просит предоставить сводный отчет по всем пользователем, с указанием сумм выручки по каждой из предоставленной услуги для расчета и уплаты налогов.
**Задача:** реализовать метод для получения месячного отчета. На вход: год-месяц. На выходе ссылка на CSV файл.
**Пример отчета:**

    название услуги 1;общая сумма выручки за отчетный период
    название услуги 2;общая сумма выручки за отчетный период

[x] **Доп. задание 2:**
Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства.
**Задача:** необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате.

### Запуск 
1. Склонировать репозиторий на локальную машину 
    ```shell
    git clone https://github.com/fdddfdfs/avito_test.git
    cd avito_test
   ```
2. Запустить сервисы в **Docker**
    ```shell
    docker-compose up -d
    ```
3. **Swagger** документация доступна по адресу (по умолчанию)
    ```shell
    localhost:3200
    ```

    
