package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	ch_brasilapi := make(chan *http.Response)
	ch_viacep := make(chan *http.Response)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		//req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/1", nil)
		req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/01007040/json/", nil)
		if err != nil {
			panic(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		//time.Sleep(200 * time.Millisecond)
		ch_viacep <- res
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		//req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/1", nil)
		req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/01153000", nil)
		if err != nil {
			panic(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		//time.Sleep(200 * time.Millisecond)
		ch_brasilapi <- res
	}()

	select {
	case resposta := <-ch_brasilapi:
		fmt.Println("resposta 'brasilapi.com.br' :")
		io.Copy(os.Stdout, resposta.Body)
		if resposta.StatusCode != http.StatusOK {
			panic("Error: " + resposta.Status)
		}

	case resposta := <-ch_viacep:
		fmt.Println("resposta 'viacep.com.br' :")
		io.Copy(os.Stdout, resposta.Body)
		if resposta.StatusCode != http.StatusOK {
			panic("Error: " + resposta.Status)
		}
	}

	//time.Sleep(time.Second * 3)
}
