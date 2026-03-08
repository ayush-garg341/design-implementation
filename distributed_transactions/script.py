import asyncio
import httpx

URL = "http://localhost:8005/order"


async def send_request(client, request_id):

    response = await client.post(URL, json={})

    print(
        f"Request {request_id}: status={response.status_code}, response={response.text}"
    )


async def main():

    async with httpx.AsyncClient() as client:

        tasks = [send_request(client, i) for i in range(10)]

        await asyncio.gather(*tasks)


if __name__ == "__main__":
    asyncio.run(main())
