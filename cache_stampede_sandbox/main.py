from fastapi import FastAPI, Request
from redis.asyncio import Redis

redis = Redis(host="127.0.0.1", port=6379, decode_responses=True)

app = FastAPI()


@app.get("/hello/{hello_id}")
async def read_root(hello_id: int, request: Request):
    value = await redis.get("celebrity")
    return {
        "method": request.method,
        "url": str(request.url),
        "headers": dict(request.headers),
        "client": request.client.host,
        "from": f"From {hello_id}",
        "value": value,
    }
