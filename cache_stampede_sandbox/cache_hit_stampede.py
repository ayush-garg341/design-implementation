import asyncio
from redis.asyncio import Redis


class DistributedNoCoalescingCache:
    def __init__(self):
        self.redis = Redis(host="127.0.0.1", db=0, port=6379, decode_responses=True)
        self.local_inflight = {}

    async def get(self, key):
        # Local coalescing first (fast path)
        if key in self.local_inflight:
            return self.local_inflight[key]

        # Double-check cache (someone might have built it)
        value = await self.redis.get(key)
        if value is not None:
            return value

        # Actually rebuild
        print(f"Server rebuilding {key}")
        value = await self.expensive_rebuild(key)

        # Store in cache with TTL
        await self.redis.set(key, value)
        self.local_inflight[key] = value
        return value

    async def expensive_rebuild(self, key):
        """
        Your expensive operation (DB query, API calls, etc.)
        """
        await self.redis.incr("expensive_rebuild")
        await asyncio.sleep(1)
        return f"rebuilt_value_for_{key}"

    async def get_request_count(self):
        return {
            "expensive_rebuild": await self.redis.get("expensive_rebuild"),
        }

    async def refresh_cache(self):
        await self.redis.flushall()
