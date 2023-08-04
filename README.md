# VGANG Project

<p>
    <a href="https://redis.io/" target="_blank">
        <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white">
    </a>
    <a href="https://redis.io/" target="_blank">
        <img src="https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white">
    </a>
    <a href="https://docker.com/" target="_blank">
        <img src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white">
    </a>
</p>

## Introduction

The whole structure is based on clean architecture packaging is by technology approach. Although the preferable one is
packaging by
feature which is more scalable.
I do not use any web framework but rather used some libraries like `gurilla mux` to be speed up developing process. This
application uses Redis as a main database but due to its abstract implementation we can freely choose another technology
instead.
As a matter of Clean Code & Architecture, OOP, SOLID principles there are many situations in code that can be revised
but for the sake of time limitation deferred to the future.

## Usage

### Docker

To start project by docker in the root of project simple run:

```
docker compose up
```

### Make

You can run the application by make file but before make sure that your `redis-server` is running
,and you have installed `go` on your machine:

```
make start
```

## Api Guide

### Get All Products (Get method)

```
http://localhost:8080/product/all
```

```json
[
  {
    "id": 141809,
    "title": "ZRED Origin embroidered Shirt - burgundy - Damen",
    "sellerName": "ZRED",
    "sellerCurrency": "EUR",
    "sellerID": 714,
    "minPrice": 18.8435,
    "maxPrice": 18.8435,
    "minRetailPrice": 28.99,
    "maxRetailPrice": 28.99,
    "stock": 2,
    "link": "http://localhost:8080/p3Kaa",
    "hash": "p3Kaa"
  },
  {
    "id": 141678,
    "title": "Back to classic Shirt - linen/black - Herren",
    "sellerName": "ZRED",
    "sellerCurrency": "EUR",
    "sellerID": 714,
    "minPrice": 12.935,
    "maxPrice": 12.935,
    "minRetailPrice": 19.9,
    "maxRetailPrice": 19.9,
    "stock": 2,
    "link": "http://localhost:8080/i1Kaa",
    "hash": "i1Kaa"
  }
]
```

### Get Product By ShortLink (Get method)

```
http://localhost:8080/p3Kaa
```

```json
{
  "id": 141809,
  "title": "ZRED Origin embroidered Shirt - burgundy - Damen",
  "sellerName": "ZRED",
  "sellerCurrency": "EUR",
  "sellerID": 714,
  "minPrice": 18.8435,
  "maxPrice": 18.8435,
  "minRetailPrice": 28.99,
  "maxRetailPrice": 28.99,
  "stock": 2,
  "link": "http://localhost:8080/p3Kaa",
  "hash": "p3Kaa"
}
```

## Functionality

When the program starts based on what config you choose, the collector send requests to fetch data fro the remote host.
There are two methods:

1. One request to fetch all data ( could be just one time or done in an interval )
2. Concurrent multiple requests to fetch data based on pagination ( could be just one time or done in an interval )

Based on `Collector` config, you can select one of these strategies, It is fully configurable.

```yaml
Collector:
  Concurrent: true
  Interval: 30 # 0 means once at time, and 0 > means an interval in second
  SplitFactor: 30 # Be used in pagination to do collecting in concurrent
```

For instance this config get all data by multiple go routines and paginating total requests for count 30 and repeat the
collecting process every 30 seconds. Each time the list gets updated.

#### Note:

Due to simplicity I just store some data of products. It is easy to add the whole information afterward.

## Architecture

Based of clean architecture I divide the whole application into tree main sections:

* app
* domain
* infra

For the sake of decoupling, all tech related components goes in to `infra` and the business logic relies on `domain`
section. the `app` section also is responsible for swing both `domain` and `infra` layers.

## Components

I decided to add some extra components to deliver high quality application. These components are:

#### Fanin Concurrency Pattern

By this pattern I handled the concurrent way of sending request to the host to fetch data. I

#### Retry Pattern

On connecting to the redis for instance, application seeks the connection process for some times fulfill the job. But
this could be used also in database queries.

#### Auto Healer Mechanising

In an cloud native environment, applications should keep themselves up as could as possible. This cloud be one of
components to keep infra connections heal themselves if the connection is lost in runtime by creating new connection and
substitute with prior one.

#### RateLimiter

Its functionality is obvious. It sets in the middle of all processes and check if the user exceeds its quota or not.

## Drawbacks

I've tried to deliver a high quality code. But still there are some drawback that due to the limited time I decided
to put it off. One of them was not handling well the shut-downing process.
Although the code base is fully testable, I didn't write **tests** for main components. This would be my next step to
fully
complete the task.