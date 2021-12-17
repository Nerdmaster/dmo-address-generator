Make addresses for your DMO to confuse and annoy others! It's fun!

- Install Go
- Clone the project
- Run "make" or simply "go build"
- ???
- Profit!

Example command:

```bash
# Generate 10,000 addresses into a wallet called "Mining", each with a
# sequential label such as "Mining 00000001"
./genaddrs http://192.168.0.8:6433 user 123456 Mining "Mining #%08d" 10000
```

Then generate an address list for your miner to rotate through:

```bash
curl --user doggles:9988 --data-binary \
  '{"jsonrpc": "1.0", "id": "curltest", "method": "listreceivedbyaddress", "params": [0, true]}' -H 'content-type: text/plain;' \
  http://192.168.0.8:6433/wallet/Mining | \
      jq '.result[] | select (.label | contains ("Mining")) .address' >> addrs.txt
```
