/*
This a RES service that shows a lists of cards, where card prices can be added,
edited and deleted by multiple users simultaneously.

* It resets the collection and models on server restart.
* It serves a web client at http://localhost:8083
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jirenius/go-res"
)

// card represents a card model
type card struct {
	ID     int64  `json:"id"`
	Price  string `json:"price"`
	Style string `json:"style"`
	PrevPrice  string `json:"prevprice"`
	Signal string `json:"signal"`
	TradeBuy string `json:"tradebuy"`
	TradeSell string `json:"tradesell"`
	Instrument string `json:"instrument"`
}

// Map of all card models
var cardModels = map[string]*card{
	"stock.card.1": {ID: 1, Price: "1500.13 ▲ ", Style: "h2", PrevPrice: "1500.12", Signal: "Neutral", TradeBuy: "0", TradeSell: "0", Instrument: "XAU/USD"},
	"stock.card.2": {ID: 2, Price: "1500.14 ▲ ", Style: "h2", PrevPrice: "1500.12", Signal: "Neutral", TradeBuy: "0", TradeSell: "0", Instrument: "GBP/USD"},
	"stock.card.3": {ID: 3, Price: "1500.15 ▼ ", Style: "h1", PrevPrice: "1501.12", Signal: "Neutral", TradeBuy: "0", TradeSell: "0", Instrument: "EUR/USD"},
	"stock.card.4": {ID: 4, Price: "1500.15 ▼ ", Style: "h1", PrevPrice: "1501.12", Signal: "Neutral", TradeBuy: "0", TradeSell: "0", Instrument: "OIL CR"},
	"stock.card.5": {ID: 5, Price: "1500.15 ▼ ", Style: "h1", PrevPrice: "1501.12", Signal: "Neutral", TradeBuy: "0", TradeSell: "0", Instrument: "OIL BR"},
}

// Collection of cards
var cards = []res.Ref{
	res.Ref("stock.card.1"),
	res.Ref("stock.card.2"),
	res.Ref("stock.card.3"),
	res.Ref("stock.card.4"),
	res.Ref("stock.card.5"),
}

// ID counter for new card models
var nextcardID int64 = 6

func main() {
	// Create a new RES Service
	s := res.NewService("stock")

	// Add handlers for "stock.card.$id" models
	s.Handle(
		"card.$id",
		res.Access(res.AccessGranted),
		res.GetModel(getcardHandler),
		res.Set(setcardHandler),
	)

	// Add handlers for "stock.cards" collection
	s.Handle(
		"cards",
		res.Access(res.AccessGranted),
		res.GetCollection(getcardsHandler),
		res.Call("new", newcardHandler),
		res.Call("delete", deletecardHandler),
	)

	// Run a simple webserver to serve the client.
	// This is only for the purpose of making the example easier to run.
	go func() { log.Fatal(http.ListenAndServe(":8083", http.FileServer(http.Dir("wwwroot/")))) }()
	fmt.Println("Client at: http://localhost:8083/")

	s.ListenAndServe("nats://localhost:4222")
}

func getcardHandler(r res.ModelRequest) {
	card := cardModels[r.ResourceName()]
	if card == nil {
		r.NotFound()
		return
	}
	r.Model(card)
}

func setcardHandler(r res.CallRequest) {
	card := cardModels[r.ResourceName()]
	if card == nil {
		r.NotFound()
		return
	}

	// Unmarshal parameters to an anonymous struct
	var p struct {
		Price  *string `json:"price,omitempty"`
		Style  *string `json:"style,omitempty"`
		PrevPrice  *string `json:"prevprice,omitempty"`
		Signal  *string `json:"signal,omitempty"`
		TradeBuy  *string `json:"tradebuy,omitempty"`
		TradeSell  *string `json:"tradesell,omitempty"`
		Instrument *string `json:"instrument,omitempty"`
	}
	r.ParseParams(&p)

	// Validate price param
	if p.Price != nil {
		*p.Price = strings.TrimSpace(*p.Price)
		if *p.Price == "" {
			r.InvalidParams("Price must not be empty")
			return
		}
	}

	// Validate style param
	if p.Style != nil {
		*p.Style = strings.TrimSpace(*p.Style)
		if *p.Style== "" {
			r.InvalidParams("Style must not be empty")
			return
		}
	}

	// Validate price param
	if p.PrevPrice != nil {
		*p.PrevPrice = strings.TrimSpace(*p.PrevPrice)
		if *p.PrevPrice == "" {
			r.InvalidParams("PrevPrice must not be empty")
			return
		}
	}

	// Validate signal param
	if p.Signal != nil {
		*p.Signal = strings.TrimSpace(*p.Signal)
		if *p.Signal == "" {
			r.InvalidParams("Signal must not be empty")
			return
		}
	}

	// Validate tradebuy param
	if p.TradeBuy != nil {
		*p.TradeBuy = strings.TrimSpace(*p.TradeBuy)
		if *p.TradeBuy == "" {
			r.InvalidParams("Tradebuy must not be empty")
			return
		}
	}

	// Validate tradesell param
	if p.TradeSell != nil {
		*p.TradeSell = strings.TrimSpace(*p.TradeSell)
		if *p.TradeSell == "" {
			r.InvalidParams("TradeSell must not be empty")
			return
		}
	}

	// Validate instrument param
	if p.Instrument != nil {
		*p.Instrument = strings.TrimSpace(*p.Instrument)
		if *p.Instrument == "" {
			r.InvalidParams("Instrument must not be empty")
			return
		}
	}

	changed := make(map[string]interface{}, 7)

	// Check if the price property was changed
	if p.Price != nil && *p.Price != card.Price {
		// Update the model.
		card.Price = *p.Price
		changed["price"] = card.Price
	}

	// Check if the style property was changed
	if p.Style != nil && *p.Style != card.Style {
		// Update the model.
		card.Style = *p.Style
		changed["style"] = card.Style
	}

	// Check if the prevprice property was changed
	if p.PrevPrice != nil && *p.PrevPrice != card.PrevPrice {
		// Update the model.
		card.PrevPrice = *p.PrevPrice
		changed["prevprice"] = card.PrevPrice
	}

	// Check if the signal property was changed
	if p.Signal != nil && *p.Signal != card.Signal {
		// Update the model.
		card.Signal = *p.Signal
		changed["signal"] = card.Signal
	}

	// Check if the tradebuy property was changed
	if p.TradeBuy != nil && *p.TradeBuy != card.TradeBuy {
		// Update the model.
		card.TradeBuy = *p.TradeBuy
		changed["tradebuy"] = card.TradeBuy
	}

	// Check if the tradesell property was changed
	if p.TradeSell != nil && *p.TradeSell != card.TradeSell {
		// Update the model.
		card.TradeSell = *p.TradeSell
		changed["tradesell"] = card.TradeSell
	}

	// Check if the instrument property was changed
	if p.Instrument != nil && *p.Instrument != card.Instrument {
		// Update the model.
		card.Instrument = *p.Instrument
		changed["instrument"] = card.Instrument
	}

	// Send a change event with updated fields
	r.ChangeEvent(changed)

	// Send success response
	r.OK(nil)
}

func getcardsHandler(r res.CollectionRequest) {
	r.Collection(cards)
}

func newcardHandler(r res.CallRequest) {
	var p struct {
		Price  string `json:"price"`
		Style  string `json:"style"`
		PrevPrice  string `json:"prevprice"`
		Signal  string `json:"signal"`
		TradeBuy  string `json:"tradebuy"`
		TradeSell  string `json:"tradesell"`
		Instrument string `json:"instrument"`
	}
	r.ParseParams(&p)

	// Trim whitespace
	price := strings.TrimSpace(p.Price)
	style := strings.TrimSpace(p.Style)
	prevprice := strings.TrimSpace(p.PrevPrice)
	signal := strings.TrimSpace(p.Signal)
	tradebuy := strings.TrimSpace(p.TradeBuy)
	tradesell := strings.TrimSpace(p.TradeSell)
	instrument := strings.TrimSpace(p.Instrument)

	// Check if we received both price and instrument
	if price == "" || style == "" || prevprice == "" || signal == "" || tradebuy == "" || tradesell == "" || instrument == "" {
		r.InvalidParams("Must provide price, style, prevprice, signal, tradebuy, tradesell and instrument")
		return
	}
	// Create a new card model
	rid := fmt.Sprintf("stock.card.%d", nextcardID)
	card := &card{ID: nextcardID, Price: price, Style: style, PrevPrice: prevprice,
		Signal: signal, TradeBuy: tradebuy, TradeSell: tradesell, Instrument: instrument}
	nextcardID++
	cardModels[rid] = card
	cardModels[rid] = card
	cardModels[rid] = card
	cardModels[rid] = card
	cardModels[rid] = card
	cardModels[rid] = card
	cardModels[rid] = card

	// Convert resource ID to a resource reference
	ref := res.Ref(rid)
	// Send add event
	r.AddEvent(ref, len(cards))
	// Appends the card reference to the collection
	cards = append(cards, ref)

	// Respond with a reference to the newly created card model
	r.Resource(rid)
}

func deletecardHandler(r res.CallRequest) {
	// Unmarshal parameters to an anonymous struct
	var p struct {
		ID int64 `json:"id,omitempty"`
	}
	r.ParseParams(&p)

	rname := fmt.Sprintf("stock.card.%d", p.ID)

	// Ddelete card if it exist
	if _, ok := cardModels[rname]; ok {
		delete(cardModels, rname)
		// Find the card in cards collection, and remove it
		for i, rid := range cards {
			if rid == res.Ref(rname) {
				// Remove it from slice
				cards = append(cards[:i], cards[i+1:]...)
				// Send remove event
				r.RemoveEvent(i)

				break
			}
		}
	}

	// Send success response. It is up to the service to define if a delete
	// should be idempotent or not. In this case we send success regardless
	// if the card existed or not, making it idempotent.
	r.OK(nil)
}
