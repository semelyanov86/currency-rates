
# Currency rates script

Script for showing currency lists for conky app.

Example usage:

1. Create crontab entry `*/15 * * * * ~/converter > ~/curren.log`
2. Set environment variable: `CURRENCY_API=`. To get converter API key, visit website: www.currencyconverterapi.com
3. Add rule for conky:
```bash
${font Entopia:bold:size=12}${color 34a8eb}CURRENCY RATES ${hr 2}$color
${offset 0}${font Noto Sans:size=10}${execi 500 cat /data/sergeyem/currency-rates/curren.log}
```

By default, parent currency is set tu RUB and target currencies is set to EUR and USD. You can change this behavior, by using arguments:

* `parent` - use code for parent currency, for example 'CHF'
* `targets` - list of target currencies, for example 'USD,EUR'


## Installation

Download and run `converter` file
