package main

import (
	"net/http"
	"encoding/json"
	"mappi/service"
	"time"
	"fmt"
)

type Foo struct {
	Name string `json:"name"`
}

func main() {
	mux := http.NewServeMux()

	out := make(chan []byte)
	markerService := service.NewHub(out)

	go func() {
		i := 0.0
		for {
			i = i + 1.0
			lat := 58.0+i*0.1
			out <- []byte(fmt.Sprintf(`
<mp-marker
	id='xmark'
	lat='%f'
	lng='18.5'
>
<sl-icon name='x-circle'></sl-icon>
</mp-marker>`, lat))
			time.Sleep(time.Second)
		}
	}()

	mux.HandleFunc("GET /markers", func(w http.ResponseWriter, r * http.Request) {
		markerService.Subscribe(w, r)
	})

	mux.HandleFunc("GET /foo", func (w http.ResponseWriter, r * http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Foo{Name: "erik"})
	})

	mux.HandleFunc("GET /div", func (w http.ResponseWriter, r * http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div>
	<p> some content </p>
</div>
<p hx-swap-oob='true' id='stuff'> hehe the stuff </p>
<mp-marker lat='60' lng='19' id='xmark' hx-swap-oob='true'>
	<sl-icon name='x-circle'></sl-icon>
</mp-marker>`))
	})

	http.ListenAndServe("0.0.0.0:8000", mux)
}
