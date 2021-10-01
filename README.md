# Для запуска проекта выполните:
```bash
docker-compose up 
```

# Для проверки проекта в соседнем терминале:
```bash
curl http://localhost:8000/ -v
```
В ответе должна быть строка: `Всё работает!`

# Для того, чтобы проверить связь приложения и бд, воспользуемся другой командой:
```bash
curl -H "Content-Type: application/json" http://localhost:8000/todos/ -d '{"name":"Wash the garbage","description":"Be especially thorough"}' -v
```
В ответ прилетит json с ID нашего todo

Чтобы проверить, что у нас всё работает, из текущей папки:
1) входим в docker контейнер
```bash
docker-compose exec postgres psql -U user db
```
2) Просим показать содержание таблицы todos
```sql
select *from todos;
```
В ответе видим таблицу с нашими вводными данными:
```bash
 id |       name       |      description       
----+------------------+------------------------
  1 | Wash the garbage | Be especially thorough
(1 row)
```
