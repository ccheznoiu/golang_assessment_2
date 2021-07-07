## Objective
* A "fast and cost-efficient" **Go** API to a song releases REST service which returns individual songs by release date
* Accepts `from` and `until` URL parameter dates
* Optionally accepts `artist` URL parameter
* Groups songs by date
* The downstream serves two URI's, returning daily and monthly results respectively, with a cost-to-cost ratio of 1:25
* These results are "usable" for 30 days
## Rationale
* It is preferrable to use the monthly URI when:
    * a request is new (not cached) and spans 25 or more days of a month
    * request volume is high, i.e. had a partial (fewer than 25 days) request been previously made for a given month, requests for the missing portion (or the entirety) of the month can be expected within the 30-day lifespan of the response data
    * high request volume also dictates that the response cache may be prebuilt, i.e. by requesting all months at once as opposed to per request (which also prioritizes _speed_ with no impact to _cost-efficiency_)
* It is preferrable to use the daily URI when:
    * a request is new and partial
    * request volume is low (converse of above)
    * data volume is high, i.e. even given high volume requests, the accompanying assumptions are attenuated
    * release data is up-to-date and requests for the current month can be expected
* High volume data would dictate a presisted cache, but we will prefer a memory cache for Go demonstration purposes
* We will prefer "cost-effectiveness" since a well-written Go program can be expected to be "fast." As such, and given no indication of the two volumes described above, we will populate the cache on a per-request basis (the latency of the downstream service will still fall upon requesters of new data).
* Given no error handling requirements, expired data would serve when the downstream is unexpectedly unavailable (unless there were a contractual definition of "usable")
* Nonetheless, a housekeeping goroutine which deletes expired cache entries is included for demonstration purposes
