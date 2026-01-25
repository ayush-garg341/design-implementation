#### Simulating the case where concurrent requests reaches the database under heavy load when there is cache miss. And avoiding this as well.

- This project demonstrate that under heavy load if a hot key expires, many requests on webserver will reach db to re-build the cache.
- If the db operation is complex and time taking, many requests can choke the database for normal ( short lived ) read operations as well.
- If we have correct locking in place, then we can allow only one request to re-build the cache for that hot key and other requests will wait.
- This problem is called cache stampede.


In one terminal:
- ulimit -n 20000
- python3 -m venv venv
- source venv/bin/activate
- pip install -r requirements.txt
- python3 -m uvicorn main:app --reload

In other terminal:
- docker-compose -f redis-cluster/docker-compose.yml up -d
- source venv/bin/activate
- python3 req_hitter.py
- curl http://127.0.0.1:8000/request_count -> count the number of expensive rebuilds simulating db operations.
- curl http://127.0.0.1:8000/refresh_cache -> delete all the entries from a cache.

