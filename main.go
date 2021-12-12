package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var templateRegexp = regexp.MustCompile(`%(0\d|)d`)

func usage(message string) {
	fmt.Fprintln(os.Stderr, message)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Usage: %s <url> <username> <password> <wallet name> <label template> <number of addresses to create>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 7 {
		usage("Not enough args")
	}

	var urlString, user, pass, wallet, template, numaddrStr = os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]
	var numAddrs, _ = strconv.Atoi(numaddrStr)
	if numAddrs < 1 {
		usage(fmt.Sprintf("Invalid number of addresses value %q", numaddrStr))
	}

	var u, err = url.Parse(urlString)
	if err != nil {
		usage(fmt.Sprintf("Invalid URL %q: %s", urlString, err))
	}

	var matches = templateRegexp.FindAllStringIndex(template, -1)
	if len(matches) != 1 {
		usage(fmt.Sprintf(`Invalid template %q: template must include one instance of "%%d" or "%%0nd"`, template))
	}

	u.User = url.UserPassword(user, pass)
	u.Path = "/wallet/" + wallet

	var addr string
	for i := 0; i < numAddrs; i++ {
		var data = bytes.NewBufferString(`{ "jsonrpc":"1.0", "id":"curltest", "method":"getnewaddress", "params": ["` +
			fmt.Sprintf(template, i) + `"]}`)
		addr, err = genAddr(u, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to POST to URL %q: %s", u.String(), err)
			fmt.Fprintln(os.Stderr, "Trying again in five seconds")
			i--
			time.Sleep(time.Second * 5)
		}
		fmt.Println(addr)
		time.Sleep(time.Millisecond * 250)
	}
}

func genAddr(u *url.URL, data io.Reader) (addr string, err error) {
	var r *http.Response
	r, err = http.Post(u.String(), "text/plain", data)
	if err != nil {
		return "", err
	}

	var body []byte
	body, err = io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	addr = resp["result"].(string)
	return addr, nil
}
