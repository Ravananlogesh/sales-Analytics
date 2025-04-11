# Go API Sales Analytics

## Tech stack

- Language : Go 1.2+
- CSV : Encoding/csv
- Schedular : fix the time to run go routine

## External Package

- Gin
- Toml
- UUID

### clone URL

### **1️⃣ Install Dependencies**

```sh
go mod tidy
```

## sample Request

```sh
curl http://localhost:28090/revenue/by-product?s_date=2012-03-12^&e_date=2025-04-10

curl http://localhost:28090/revenue/total/2025-01-01/2025-04-30


curl http://localhost:28090/refersh

```

## sample Response
- Endpoint :  /revenue/by-product?s_date=2012-03-12^&e_date=2025-04-10

```json
{
  "data": {
    "end": "2025-04-10",
    "revenues": [
      { "product_name": "iPhone 15 Pro", "revenue": 3767.1 },
      { "product_name": "UltraBoost Running Shoes", "revenue": 504 },
      { "product_name": "Sony WH-1000XM5 Headphones", "revenue": 297.4915 },
      { "product_name": "Levi's 501 Jeans", "revenue": 143.976 }
    ],
    "start": "2012-03-12"
  },
  "message": "Revenue by product",
  "status": "S"
}


```

### /revenue/total/2025-01-01/2025-04-30
```json
{
    "data": {
        "end_date": "2025-04-30",
        "start_date": "2025-01-01",
        "total_amount": 0
    },
    "message": "total revenue",
    "status": "S"
}
```


