from fastapi import FastAPI

from contextlib import asynccontextmanager
from alembic import command
from alembic.config import Config
from router import api_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    try:
        alembic_cfg = Config("alembic.ini")
        command.upgrade(alembic_cfg, "head")
        yield
    except Exception as e:
        print(f"Startup error: {e}")
        raise


app = FastAPI(title="Distributed Transactions", lifespan=lifespan)
app.include_router(api_router)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
