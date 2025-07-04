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
**POST** `/api/v1/orders/process`

#### Request Body (Example)
```json
[
    {
        "no": 1,
        "platformProductId": "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
        "qty": 1,
        "unitPrice": 160,
        "totalPrice": 160
    },
    {
        "no": 2,
        "platformProductId": "FG0A-PRIVACY-IPHONE16PROMAX",
        "qty": 1,
        "unitPrice": 50,
        "totalPrice": 50
    }
]
```

#### Response
```json
{
    "data": [
        {
            "no": 1,
            "productId": "FG0A-CLEAR-OPPOA3",
            "materialId": "FG0A-CLEAR",
            "modelId": "OPPOA3",
            "qty": 2,
            "unitPrice": 40.00,
            "totalPrice": 80.00
        },
        {
            "no": 2,
            "productId": "FG0A-MATTE-OPPOA3",
            "materialId": "FG0A-MATTE",
            "modelId": "OPPOA3",
            "qty": 2,
            "unitPrice": 40.00,
            "totalPrice": 80.00
        },
        {
            "no": 3,
            "productId": "FG0A-PRIVACY-IPHONE16PROMAX",
            "materialId": "FG0A-PRIVACY",
            "modelId": "IPHONE16PROMAX",
            "qty": 1,
            "unitPrice": 50.00,
            "totalPrice": 50.00
        },
        {
            "no": 4,
            "productId": "WIPING-CLOTH",
            "qty": 5,
            "unitPrice": 0.00,
            "totalPrice": 0.00
        },
        {
            "no": 5,
            "productId": "CLEAR-CLEANNER",
            "qty": 2,
            "unitPrice": 0.00,
            "totalPrice": 0.00
        },
        {
            "no": 6,
            "productId": "MATTE-CLEANNER",
            "qty": 2,
            "unitPrice": 0.00,
            "totalPrice": 0.00
        },
        {
            "no": 7,
            "productId": "PRIVACY-CLEANNER",
            "qty": 1,
            "unitPrice": 0.00,
            "totalPrice": 0.00
        }
    ],
    "status": "success"
}
```

### Health Check
**GET** `/health`
