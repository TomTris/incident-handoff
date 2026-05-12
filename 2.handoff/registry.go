package main

type Registry struct {
	hubs          map[string]*Hub
	clientCounter map[string]int
	register      chan *Client
	unregister    chan *Client
	broadcast     chan BroadcastMessage
}

type BroadcastMessage struct {
	msg        string
	incidentID string
}

func NewRegistry() Registry {
	return Registry{
		hubs:          make(map[string]*Hub),
		clientCounter: make(map[string]int),
		register:      make(chan *Client), // no buffered on purpose
		unregister:    make(chan *Client), // no buffered on purpose
		broadcast:     make(chan BroadcastMessage),
	}
}

func (r *Registry) run() {
	for {
		select {
		case client := <-r.register:
			r.joinRegistry(client)
		case client := <-r.unregister:
			r.leaveRegister(client)
		}
	}
}

func (r *Registry) joinRegistry(client *Client) {
	incID := client.incidentID
	r.clientCounter[incID]++

	// If first Client
	if r.clientCounter[incID] == 1 {
		r.hubs[incID] = NewHub()
		go r.hubs[incID].run()
	}
	client.hub = r.hubs[incID]
	client.hub.register <- client
}

func (r *Registry) leaveRegister(client *Client) {
	incID := client.incidentID
	r.clientCounter[incID]--

	hub, _ := r.hubs[incID]
	if r.clientCounter[incID] == 0 {
		r.hubs[incID] = nil
		close(hub.done)
		return
	}
	hub.unregister <- client
}
