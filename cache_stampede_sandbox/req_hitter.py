import asyncio
import aiohttp
import time

URL = "http://127.0.0.1:8000/hello"

# Number of concurrent requests we want
CONCURRENCY = 10000

responses = []


async def fetch(session, i):
    """
    This function represents ONE logical request.

    - It does NOT block the thread
    - It suspends itself when waiting for network I/O
    """

    url = f"{URL}/{i}"
    async with session.get(url) as response:
        # WAIT here until response body is received
        # While waiting, the event loop runs other tasks
        t = await response.text()
        responses.append(t)
        return i


async def main():
    """
    This is the main coroutine.

    - creates the HTTP session
    - schedules thousands of tasks
    - waits for all of them to finish
    """

    # Set global timeout for requests, all requests should finish in this 30 seconds
    timeout = aiohttp.ClientTimeout(total=60)

    # how many sockets can be open at once
    connector = aiohttp.TCPConnector(limit=10)

    async with aiohttp.ClientSession(timeout=timeout, connector=connector) as session:
        # Create 10k asyncio Tasks
        tasks = [asyncio.create_task(fetch(session, i)) for i in range(CONCURRENCY)]

        # WAIT until all tasks finish
        # This does NOT block the OS thread
        await asyncio.gather(*tasks)


start = time.time()
asyncio.run(main())
print(responses, len(responses))
print(f"Completed in {time.time() - start:.2f}s")

# Add request count here for total requests, requests rebuilt, hit cache
