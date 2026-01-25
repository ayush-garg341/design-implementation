import asyncio
from redis.asyncio import Redis


class DistributedCoalescingCache:
    def __init__(self):
        self.redis = Redis(host="127.0.0.1", port=6379, db=0, decode_responses=True)
        self.local_inflight = {}

    async def get(self, key):
        # Local coalescing first (fast path)
        if key in self.local_inflight:
            return self.local_inflight[key]

        lock_key = f"lock:{key}"
        # Try to acquire distributed path
        lock = self.redis.lock(lock_key, timeout=0.5, blocking=False)

        acquired = await lock.acquire(blocking=False)
        if acquired:
            # I won the lock - I'll fetch
            try:
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
            finally:
                # Always release lock
                await lock.release()
        else:
            # Someone else has the lock, wait for them to finish
            print(f"Server waiting for {key}")
            return await self.wait_for_rebuild(lock_key, key)

    async def wait_for_rebuild(self, lock_key, key):
        """
        Wait for another server to complete the rebuild
        """
        await self.redis.incr("wait_for_rebuild")
        while True:
            value = await self.redis.get(key)
            if value is not None:
                print(f"Server got rebuilt value for {key}")
                return value

            lock_exists = await asyncio.to_thread(self.redis.exists, lock_key)
            if not lock_exists:
                print(f"Lock released but cache empty, retrying {key}")
                return await self.get(key)

    async def expensive_rebuild(self, key):
        """
        Your expensive operation (DB query, API calls, etc.)
        """
        await self.redis.incr("expensive_rebuild")
        await asyncio.sleep(0.25)
        return f"rebuilt_value_for_{key}"

    async def get_request_count(self):
        return {
            "expensive_rebuild": await self.redis.get("expensive_rebuild"),
            "wait_for_rebuild": await self.redis.get("wait_for_rebuild"),
        }

    async def refresh_cache(self):
        await self.redis.flushall()
