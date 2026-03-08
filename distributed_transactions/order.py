import httpx
import uuid
from fastapi import FastAPI, HTTPException


app = FastAPI(title="Order Service")


@app.post("/order")
async def create_order():
    order_id = str(uuid.uuid4())
    print("order id", order_id)
    async with httpx.AsyncClient() as client:
        response = await client.post(
            "http://localhost:8000/api/food/reserve", json={"order_id": order_id}
        )

        if response.status_code != 200:
            raise HTTPException(status_code=400, detail=response.text)

        response = await client.post(
            "http://localhost:8000/api/agent/reserve", json={"order_id": order_id}
        )
        if response.status_code != 200:
            raise HTTPException(status_code=400, detail=response.text)

        response = await client.post(
            "http://localhost:8000/api/food/book", json={"order_id": order_id}
        )
        if response.status_code != 200:
            raise HTTPException(status_code=400, detail=response.text)

        response = await client.post(
            "http://localhost:8000/api/agent/book", json={"order_id": order_id}
        )
        if response.status_code != 200:
            raise HTTPException(status_code=400, detail=response.text)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8005)
