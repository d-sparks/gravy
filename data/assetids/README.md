# Tokenize

## Generate

1. Read all distinct (asset, exchange) pairs from dailyprices
2. Find the maximum assigned id (or 1)
3. For each (asset, exchange) pair, if there is no id:
    1. Insert (asset, exchange, id) into db
    2. Increase id by 1

Want unique indexes on id and on (asset, exchange).

## Consume

Since things were historically done with ticker strings (and ignoring exchange perhaps), things will probably continue to work with strings for a while. The daily prices DB can be extended to offer the mapping. Processes have access to dailyprices DB from the registrar, so each process can go back and forth between ids and (ticker, exchange) pairs.

Eventually, should try to eliminate the use of ticker strings in almost all processes.