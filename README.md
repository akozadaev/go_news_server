# Постановка задачи

Написать json REST сервер с двумя ручками:
- `POST /edit/:Id` - изменение новости по Id
- `GET /list` - список новостей

БД - mysql (можно и postgres, мы перешли на него).

В качестве сервера использовать fiber. Для работы с базой reform.

Соединение с базой должно использовать connection pool. Все настройки через переменные окружения и/или viper.

БД:
```SQL
--
-- Структура таблицы `News`
--

create table news
(
    id      bigserial not null primary key,
    title   text      not null,
    content text      not null
);
-- --------------------------------------------------------

--
-- Структура таблицы `NewsCategories`
--

create table newscategories
(
    category_id bigserial not null primary key,
    news_id     bigint    not null,
    foreign key (news_id) references news
);
```
Т.е. легко видеть, что связка новостей и категорий идёт в отдельной таблице.

Формат входных данных для первой ручки:

```json
{
  "Id": 64,
  "Title": "Lorem ipsum",
  "Content": "Dolor sit amet <b>foo</b>",
  "Categories": [1,2,3]
}
```

При этом, если какое-то из полей не задано - это поле не должно быть обновлено.

Формат данных на выходе list:

```json
{
    "Success": true,
    "News": [
      {
        "Id": 64,
        "Title": "Lorem ipsum",
        "Content": "Dolor sit amet <b>foo</b>",
        "Categories": [1,2,3]
      },
      {
        "Id": 1,
        "Title": "first",
        "Content": "tratata",
        "Categories": [1]
      }
    ]
}
```

## Требования и пожелания

1. Срок выполнения - не больше трёх дней. Мы высылаем тестовые задания пачками сразу десяткам людей, поэтому вас просто могут опередить другие кандидаты.
2. Если знакомы с docker - хотелось бы посмотреть упаковку сервиса в контейнер.
3. Дополнительно хотелось бы увидеть(плюс в карму по сравнению с другими кандидатами):
- авторизацию через Authorization заголовок и грамотную структуризацию кода и ручек по группам/папкам;
- валидацию полей при редактировании;
- пагинацию на ручке списка;
- грамотное логгирование с использованием любого популярного логгера(напр. logrus);
- грамотную обработку ошибок.

______________
# Реализация

Реализованы 4 маршрута. Два без авторизации:
[GET] localhost:8181/list?limit=1&offset=0

ответ:

```{
    "news": [
        {
            "id": 1,
            "title": "title1",
            "content": "content2",
            "categories": null
        }
    ],
    "success": true
}
```
[POST] localhost:8181/edit/1

body:

```
{
    "id": 1,
    "Title": "title1",
    "Content": "content2",
    "categories": [1,2]
}
```
ответ:

```
{
    "success": true
}
```

Два с авторизацией:
[GET] localhost:8181/private/list?limit=1&offset=0
[POST] localhost:8181/private/edit/1

Логгер - zap
БД - PostgeSQL 
