
# Currency rates script

Script for showing currency lists for conky app.

Example usage:

1. Create crontab entry `*/15 * * * * ~/converter > ~/curren.log`
2. Add rule for conky:
```bash
${font Entopia:bold:size=12}${color 34a8eb}CURRENCY RATES ${hr 2}$color
${offset 0}${font Noto Sans:size=10}${execi 500 cat /data/sergeyem/currency-rates/curren.log}
```


## Installation

Download and run `converter` file
