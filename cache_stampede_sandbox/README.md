#### Simulating the case where concurrent requests reaches the database under heavy load when there is cache miss. And avoiding this as well.

In one terminal:
- ulimit -n 20000
- python3 -m venv venv
- source venv/bin/activate
- pip install -r requirements.txt
- uvicorn main:app --reload

In other terminal:
- docker-compose -f redis-cluster/docker-compose.yml up -d
- source venv/bin/activate
- python req_hitter.py

- Simulating Database i.e some complex operation by file reading doing some complex op.
