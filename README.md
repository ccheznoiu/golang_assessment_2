## Build & Run
* `cd ccheznoiu-dd-technical-test`
* `export BV_APIKEY='' # supply secret key`
* `go build`
* `./ccheznoiu-dd-technical-test &`
* `curl http://localhost:8000/releases?from=2021-01-01&until=2021-03-31 # etc.`

## Rationale
###### Cost-efficiency is preferred over performance.
Because a `/monthly` call is cheaper than a `/daily` call for **every** day of a given month, prefer the former when `from` and `until` span 25 days of a month or more.

An in-memory cache stores already-requested dates for 30 days. It is populated **per request** which puts the associated latency on first-time requesters. A performance-preferential design would, e.g. query every month first to build the cache, making the daily endpoint unnecessary (unless the downstream data is updated daily).

Implied is **high-volume data**, i.e. making it unfeasible to keep a full copy of the datastore in memory in addition to making requests unpredictable.

Also implied is **low-volume traffic**, i.e. we do not assume that requests for any given date can land on any given day, which would have dictated that a missing month would eventually be filled out by partial requests, making it cheaper to fill out the entire month once it was discovered to be missing.

## Future Enhancements
* consider a persisted cache for case: _request over unfeasible range of dates_
* pass error response from Release Service to response body
