# desafio-go-bases

## How to run the app

For running the app you just need a csv with the following information in each row: id, name, email,
destiny, flight time, ticket cost, with no headers

Then you just need to copy the following command

```shell
go run main.go --input tickets.csv --destination China --total 500
```
Parameters:
- input -> the file to be processed. `no default value`
- destination -> which country do you want to evaluate `default: Brazil`
- total -> the amount of people to compare with `default: 1000`