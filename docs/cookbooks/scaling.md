# Horizontal Scaling

Since *simple-auth* uses signed cookies to manage the session, horizontal scaling is easy.  Simply running more than one instance side-by-side leads to easy scaling.

::: warning
To most effectively scale, we recommend using postgresql or mysql/mariadb. You can read more about using different databases [here](../database.md)
:::

You need to make sure *simple-auth* is behind some sort of load balancer.  This example uses [traefik](./traefik.md), but nginx or swarm balancing should work fine.


## docker-compose

Creating replicas is as simple as adding this in the service:

<<< @/examples/docker-scaling/docker-compose.yml#replicas

Full `docker-compose.yml`

<<< @/examples/docker-scaling/docker-compose.yml
