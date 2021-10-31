# Crypto Balances Exporter
exposes your balances and initial costs as prometheus metrics 

## Usage
create a config and pass it to the exporter
````
coins:
  - name: ETH
    amount: 100.039996
  - name: BTC
    amount: 53.8071
    totalCost: 100 #optional

````

### Parameters
```
Usage of crypto-balances:
  -config string
        Path to configfile (default "config.yaml")
  -debug
        Set debug log level.
  -httpServerPort uint
        HTTP server port. (default 9101)

```
