const {Elem, Txt} = window["modapp-base-component"];
const {CollectionList, ModelTxt} = window["modapp-resource-component"];
const ResClient = resclient.default;

// Creating the client instance.
let client = new ResClient('ws://localhost:8080');

// Error handling
let errMsg = new Txt();
let errTimer = null;
errMsg.render(document.getElementById('error-msg'));
let showError = (err) => {
    errMsg.setText(err && err.message ? err.message : String(err));
    clearTimeout(errTimer);
    errTimer = setTimeout(() => errMsg.setText(''), 7000);
};

// Get the collection from the service.
client.get('stock.cards').then(cards => {
    // Render the collection of cards
    new Elem(n =>
        n.component(new CollectionList(cards, card => {
            let c = new Elem(n =>
                n.elem('div', {className: 'list-item'}, [
                    n.elem('div', {className: 'card shadow'}, [
                        // View card mode
                        n.elem('div', {className: 'view'}, [
                            n.elem('div', {className: 'instrument'}, [
                                n.component(new ModelTxt(card, card => card.instrument, {tagName: 'h3'}))
                            ]),
                            n.elem('div', {className: 'price'}, [
                                n.component(new ModelTxt(card, card => card.price, {tagName: card.style}))
                            ]),
                            n.elem('div', {className: 'instrument'}, [
                                n.component(new Txt("(")),
                                n.component(new ModelTxt(card, card => card.prevprice)),
                                n.component(new Txt(")"))
                            ]),
                            n.elem('div', {className: 'instrument'}, [
                                n.component(new Txt("Signal: ")),
                                n.component(new ModelTxt(card, card => card.signal))
                            ]),
                            n.elem('div', {className: 'instrument'}, [
                                n.component(new Txt("Buy Trades: ")),
                                n.component(new ModelTxt(card, card => card.tradebuy))
                            ]),
                            n.elem('div', {className: 'instrument'}, [
                                n.component(new Txt("Sell Trades: ")),
                                n.component(new ModelTxt(card, card => card.tradesell))
                            ])
                        ])
                    ])
                ])
            );
            return c;
        }, {className: 'list'}))
    ).render(document.getElementById('stock'));
}).catch(err => showError(err.code === 'system.connectionError'
    ? "Connection error. Are NATS Server and Resgate running?"
    : err
));
