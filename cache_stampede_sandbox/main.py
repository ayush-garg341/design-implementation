from fastapi import FastAPI, Request

# from cache_avoid_stampede import DistributedCoalescingCache

from cache_hit_stampede import DistributedNoCoalescingCache

app = FastAPI()

# redis_helper = DistributedCoalescingCache()

redis_helper = DistributedNoCoalescingCache()


@app.get("/hello/{hello_id}")
async def read_root(hello_id: int, request: Request):
    value = await redis_helper.get("celebrity")
    return {
        "method": request.method,
        "url": str(request.url),
        "headers": dict(request.headers),
        "client": request.client.host,
        "from": f"From {hello_id}",
        "value": value,
    }


@app.get("/request_count")
async def request_count():
    return {"request_count": await redis_helper.get_request_count()}


# Deletes all the keys from cache, useful when running multiple and different iterations
@app.get("/refresh_cache")
async def refresh_cache():
    await redis_helper.refresh_cache()
    return {"success": "done"}
