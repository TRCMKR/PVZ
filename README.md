## Утилита для управления ПВЗ

### CLI Команды

- `help` – выводит справку
- `clear` – очищает терминал
- `acceptOrder <oid> <uid> <wght> <prc> <expd> <pkg1?> <pkg2?>` – принимает заказ
- `acceptOrders <path>` – принимает заказы в json-файле
- `returnOrder <oid>` – возвращает заказ
- `processOrders <uid> ...<oids> <action>` – выдаёт заказы пользователю/принимает возвраты от пользователя
- `userOrders <uid>` – выводит заказы пользователя
- `returns` – выводит все возвращённые заказы
- `orderHistory` – выводит историю всех заказов
- `exit` – сохраняет данные и закрывает приложение

### Web Команды
- `/orders [get]` – получает список заказов, фильтрация на все поля, кроме id и даты последнего изменения
```bash
curl -u lol:12345678 --request GET \
"localhost:9000/orders?user_id=789&count=5&page=1"
```
- `/orders [post]` – создаёт новый заказ
```bash
curl -u lol:12345678 --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1009,"user_id":52,"weight":100,"price":{"amount":1000000,"currency":"RUB"},"packaging":"box","extra_packaging":"wrap","expiry_date":"2025-03-10T00:00:00Z"}' \
"http://localhost:9000/orders"
```
- `/orders/{id} [delete]` – удаляет заказ
```bash
curl -u lol:12345678 --request DELETE \
"http://localhost:9000/orders/1009"
```
- `/orders/process [post]` – обрабатывает заказы пользователя
```bash
curl -u lol:12345678 --header "Content-Type: application/json" \
--request POST \
--data '{"user_id":789,"order_ids":[1009],"action":"give"}' \
http://localhost:9000/orders/process
```

- `/admins [post]` – создаёт админа
```bash
curl --header "Content-Type: application/json" \
--request POST \
--data '{"id":2,"username":"lol","password":"12345678"}' \
http://localhost:9000/admins
```
- `/admins/{username} [post]` – обновляет пароль админа
```bash
curl --header "Content-Type: application/json" \
--request POST \
--data '{"password":"12345678","new_password":"5555"}' \
http://localhost:9000/admins/lol
```
- `/admins/{username} [delete]` – удаляет админа
```bash
curl --header "Content-Type: application/json" \
--request DELETE \
--data '{"password":"5555"}' \
http://localhost:9000/admins/lol
```

### Запуск

`make build && make run` – собирает приложение и запускает

`make build-windows && make run` – собирает для винды и запускает

### Помощь по makefile

`make help` – выводит справку по всем make-таргетам