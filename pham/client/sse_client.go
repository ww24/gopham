package client

import (
	"bufio"
	"log"
	"net/http"
	"strings"

	"github.com/ww24/gopham/pham"
)

// SSEClient is Server-Sent Events client
type SSEClient struct {
	http.CloseNotifier
	res      *http.Response
	Listener <-chan pham.JSON
}

// CloseNotify implemente http.Closenotifier
func (client *SSEClient) CloseNotify() (ch <-chan bool) {
	if notifier, ok := client.res.Body.(http.CloseNotifier); ok {
		ch = notifier.CloseNotify()
		log.Println("close notify")
	}
	return
}

// Close method is alias of res.Body.Close method
func (client *SSEClient) Close() (err error) {
	err = client.res.Body.Close()
	return
}

// NewSSEClient is SSEClient constructor
func NewSSEClient(url string) (client *SSEClient, err error) {
	client = new(SSEClient)
	listener := make(chan pham.JSON)
	client.Listener = listener

	// connect to server
	res, err := http.Get(url)
	if err != nil {
		return
	}
	client.res = res

	go func() {
		defer res.Body.Close()

		sc := bufio.NewScanner(res.Body)
		sc.Split(split())

		for {
			switch {
			case res.Close:
				log.Println("closed")
				return
			case sc.Scan():
				lines := strings.Split(sc.Text(), "\n")

				data := make(pham.JSON)

				property := ""
				for _, line := range lines {
					splits := strings.SplitN(line, ":", 2)
					switch strings.Trim(splits[0], " ") {
					case "event":
						property = "event"
						data["event"] = splits[1]
					case "data":
						property = "data"
						data["data"] = splits[1]
					default:
						if str, ok := data[property].(string); ok {
							data[property] = str + "\n" + line
						}
					}
				}

				listener <- data
			}
		}
	}()

	return
}
