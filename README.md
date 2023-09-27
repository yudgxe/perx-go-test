# perx-go-test

### Run
go test ./...  
go run .\cmd\api\main.go  
```
-addr string
        Адрес для запускаемого сервера (default ":8080")
  -q int
        Количество воркеров (default 2)
  -s int
        Размер очереди (default 10)
  -v int
        Уровень «многословности» логов (default 1)
```

### docs
[Сылка](https://app.getpostman.com/join-team?invite_code=5000a014b854b70a83a2cd686349b3cb&target_code=4bbae5309fbdf34772ee8f1d55eee61c) на постман с готовыми запросами.  

Имеется два рута:  
 * /api/task/create [post]
 * /api/task/list [get]

Пример данных для рута `/api/task/create`
```
{
    "n" : 10,
    "d" : 10,
    "n1" : 10,
    "l" : 10,
    "ttl" : 10
}
```
Так же рут `/api/task/list` имеет опциональный `boolean` параметр `sorted`, указывающий нужна ли сортировка или нет, cортировка выполняется по task id.  
 
Пример рута `/api/task/list?sorted=true`

Пример овевта
```
[
    {
        "id": 1,
        "n": 10,
        "d": 10,
        "n1": 10,
        "l": 10,
        "ttl": 10,
        "iteration": 0,
        "created_at": "2023-09-27T21:22:12.6237337+05:00",
        "started_at": "2023-09-27T21:22:12.6237337+05:00",
        "ended_at": null,
        "status": "processing",
        "result": 0,
        "queue_number": 0
    }
]
```
Поле `result` хранит промежуточный/конечный результат.

### TODO
 * graceful shutdown для воркеров.
