# MeliProxy

It is a proxy that rate limit for the Mercado Libre api 

## Installation

Use docker-compose for run this.

```bash
sudo docker-compose down && sudo docker-compose build && sudo docker-compose up
```

For run more instances of the core add --scale 

```bash
sudo docker-compose down && sudo docker-compose build && sudo docker-compose up --scale meliproxy=2
```
## Usage


Example: 

http://localhost:4000/categories/MLA97994

## Architecture

Nginx balancer

Meliproxy backend GO

Redis for cache


## License
[MIT](https://choosealicense.com/licenses/mit/)