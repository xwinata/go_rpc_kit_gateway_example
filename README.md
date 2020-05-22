go run server/server.go (:8000)

go run client/client.go (:8001)

api endpoint at localhost:8001/search(POST)

body example : 
```
{
  "searchword": "avatar",
  "pagination": "1",
}
```
