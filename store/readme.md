# Store

- Places the order

## API

## Place the order

```sh
# direct
curl -X POST 0:8000/order

# via Dapr
curl -X POST 0:3500/v1.0/invoke/store/method/order
```

