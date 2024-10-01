# Warehouse

- Manages stock (initialized with 5 items)
- Notifies when out of stock


## API

### Get current stock

```sh
# direct
curl 0:8001/stock

# via Dapr
curl 0:3501/v1.0/invoke/warehouse/method/stock
```

### Decrease stock

```sh
# direct
curl -X POST -d '{"items": 1}' 0:8001/stock:decrease

# via Dapr
curl -X POST -d '{"items": 1}' 0:3501/v1.0/invoke/warehouse/method/stock:decrease
```

