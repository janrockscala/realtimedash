# Go/Rasgate/Nats Real-time Dashboard

**Tags:** *Models*, *Collections*, *Linked resources*, *Call methods*, *Resource parameters*

## Description
A simple stock monitoring system, listing instruments with prices and other indicators. All cards and fields can be edited, added, or deleted by multiple users simultaneously.

## Prerequisite

* [Download](https://golang.org/dl/) and install Go
* [Install](https://resgate.io/docs/get-started/installation/) *NATS Server* and *Resgate*

## Install and run

```text
git clone https://github.com/janrockscala/realtimedash
cd realtimedash
go run .
```

Open the client
```text
http://localhost:8083
```


## Things to try out

### Realtime updates
* Use Postman(for MacOS) or any other curl client to modify the dashboard content
* Add/Edit/Delete entries to observe realtime updates.

### System reset
* Open the client and make some changes in Go code.
* Restart the service to observe resetting of resources in all clients.

### Resynchronization
* Open the client on two separate devices.
* Disconnect one device.
* Make changes using the other device.
* Reconnect the first device to observe resynchronization.

## API

Request | Resource | Description
--- | --- | ---
*get* | `stock.cards` | Collection of card model references.
*call* | `stock.cards.new` | Creates a new card.
*get* | `stock.cards.<CARD_ID>` | Models representing card.
*call* | `stock.cards.<CARD_ID>.set` | Sets the cards' *prices* , *instruments*, *trades* etc. properties.
*call* | `stock.cards.<CARD_ID>.delete` | Deletes a card.

## REST API

Resources can be retrieved using ordinary HTTP GET requests, and methods can be called using HTTP POST requests.

### Get card collection
```
GET http://localhost:8080/api/stock/cards
```

### Get card
```
GET http://localhost:8080/api/stock/cards/<CARD_ID>
```

### Update card properties
```
POST http://localhost:8080/api/stock/card/<CARD_ID>/set
```
*Body*  
```
{ "price": "1950.34" }
```

Style: red colour text
```
{ "style": "h1" }
```
Style: red green text
```
{ "style": "h2" }
```
(custom styles are not supported)

### Add new card
```
POST http://localhost:8080/api/stock/cards/new
```
*Body*  
```
{ "instrument": "GBPUSD", "price": "1.22121" ...
```

### Delete card
```
POST http://localhost:8080/api/stock/cards/delete
```
*Body*  
```
{ "id": <CARD_ID> }
```
