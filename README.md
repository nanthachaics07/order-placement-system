# ORDER-PLACEMENT-SYSTEM
    - Golang
    - Gin
    - Unittest [coverage >= 81.3%]
    - Mockery
    - Dynamic Logger
    - Git Workflow
    (Service Processes Only NONE "context.Context" in Each func)

# Developer
 ```
 Nanthachai Intarpradit (SAFE)
 ```

## Architecture

use **Clean Architecture (Uncle Bob)**  4 layers:

```
internal
    ├── Domain Layer (Entities & Value Objects)
    ├── Use Case Layer (Business Logic)
    ├── Interface Adapters (Controllers & Presenters)
    └── Infrastructure Layer (Web Framework & External Libraries)
```


### Installation

1. **Clone the repository**
```bash
git clone https://github.com/nanthachaics07/order-placement-system.git
cd order-placement-system
```

2. **Install dependencies**
```bash
make tidy
```

3. **Run tests**
```bash
make test
```

3. **Run With test Cover**
```bash
make test-coverage
```

4. **Run the application**
```bash
make run
```

5. **Build application**
```bash
make build-up
```

6. **Down application**
```bash
make build-down
```

Server Opening on `http://localhost:8080`

##  API Endpoints

### Process Orders
**POST** `/api/orders/process`

#### Request Body (Example)
```json
{
  "orders": [
    {
      "no": 1,
      "platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
      "qty": 2,
      "unitPrice": 50.00,
      "totalPrice": 100.00
    }
  ]
}
```

#### Response
```json
{
  "success": true,
  "data": [
    {
      "no": 1,
      "productId": "FG0A-CLEAR-IPHONE16PROMAX",
      "materialId": "FG0A-CLEAR",
      "modelId": "IPHONE16PROMAX",
      "qty": 2,
      "unitPrice": 50.00,
      "totalPrice": 100.00
    },
    {
      "no": 2,
      "productId": "WIPING-CLOTH",
      "qty": 2,
      "unitPrice": 0.00,
      "totalPrice": 0.00
    },
    {
      "no": 3,
      "productId": "CLEAR-CLEANNER",
      "qty": 2,
      "unitPrice": 0.00,
      "totalPrice": 0.00
    }
  ],
  "stats": {
    "total_orders": 3,
    "main_products": 1,
    "complementary_items": 2,
    "total_quantity": 4,
    "total_value": 100.00
  }
}
```

### Health Check
**GET** `/health`
