# Demo sales app

- Store component accepts orders
- Warehouse component reduces the stock, sends events on low stock, fails on no stock
- Notifies of low stock via email

## Demonstrates

- State management
- Pubsub
- Output binding

## Run locally

- Start deps

    ```sh
    docker compose up
    ```

- Quit VPN apps (Amnesia etc), disable firewall (Lulu)

- Start all components in separate terminals using `./run.sh` scripts

- Try using `curl` to register orders with the given number of items:

    ```sh
    # Direct service invocation
    curl -X POST -H 'content-type: application/json' -d '{"items": 1}' 0:3001/order

    # Invocation via Dapr sidecar
    curl -X POST 0:3500/v1.0/invoke/store/method/order
    ```

- Try checking email at http://localhost:8025/ when remaining stock goes lower than 3.
