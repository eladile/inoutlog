
# inoutlog
log working hours to json

learn more about creating telegram bots in: https://core.telegram.org/api#bot-api

# usage
go run main.go --help 

```
Usage of inoutlog:
  -extra int
        added to each record payment (default 30)
  -out string
        where the output json is saved to (default "out.json")
  -tariff int
        the hourly tariff (default 50)
  -token string
        telegram bot token see: https://core.telegram.org/api#bot-api
```

# using the bot
```
/in [HH:MM] - will log the in entry and write the time back to the bot user.
/out - will log the out entry and write back the time + payment.
```

# records file example
```
{
  "records": [
    {
      "in": "2022-11-28T17:54:26.4788987+02:00",
      "out": "2022-11-28T18:13:47.113362+02:00",
      "totalTime": 1160634463300,
      "totalPay": 46
    },
    {
      "in": "2022-11-28T18:17:00.7759739+02:00",
      "out": "2022-11-28T18:19:33.6519995+02:00",
      "totalTime": 152878620900,
      "totalPay": 32
    },
    {
      "in": "2022-11-28T21:09:44.5137865+02:00",
      "out": "2022-11-28T22:31:05.8400719+02:00",
      "totalTime": 4881828832500,
      "totalPay": 97
    },
    {
      "in": "2022-11-28T17:00:00+02:00",
      "out": "2022-11-28T22:31:50.6291995+02:00",
      "totalTime": 19910629199500,
      "totalPay": 306
    }
  ],
  "tariff": 50,
  "extra": 30
}
```

