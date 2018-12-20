# Backend Elastic Client

## Setup
- ### Requirements
  - Elasticsearchdb
  - Imports listed in import section

```
import (
 "context"
"log"
// "encoding/json"
"fmt"
"net/http"
"math"
"reflect"
"strconv"
"time"
elastic "gopkg.in/olivere/elastic.v5"
"github.com/speps/go-hashids"
 "github.com/gin-gonic/gin"
 "github.com/gin-contrib/cors"

)```
<br>
```go get github.com/speps/go-hashids``` <br>
```github.com/gin-gonic/gin``` <br>
```github.com/gin-contrib/cors```
<br>
### Build & Run
  ```elasticsearch``` <br>
  ```go build main.go``` <br>
  ```./main``` <br>
## Endpoints
  - **Create (Generate short url in DB)**
    - url <a href> localhost:8000/create </a>
    - Data Contract
      - method - **POST**
      - (Type: application/x-www-form-urlencoded)
      - {url:""}
  - **Get short url (Get)**
    - url <a href> localhost:8000/pretty/:orig </a>
    - Data Contract
      - Method - **GET**
      - Response type **JSON** {url:""}
  - ** Get original Url**
    - URL a href> localhost:8000/redirect/:hash </a>
    - Data Contract
      - Method - **GET**
      - Response type **JSON** {url:""}
