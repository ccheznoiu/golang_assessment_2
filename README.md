## Objective
* **In Go**, build a "fast and cost-efficient" REST API to query the song releases between two dates provided as URL parameters with an optional artist filter.
* The data store is served by the _daily_ and _monthly_ URI's of a REST service, costing 1 and 25 units per each call, respectively.
* Each response is "reusable" for 30 days.
