# Simple-Auth Gateway

The below example will set up a [gateway](/access/gateway) to a backend via docker compose

## Data Flow

<mermaid>
graph LR
A{User} -- Web Request --> B[Simple-Auth]
B -- Proxies --> C[Backend]
</mermaid>

## docker-compose example

<<< @/../examples/gateway/docker-compose.yml

