import logging
from fastapi import APIRouter, Depends, HTTPException, Request
from datetime import datetime
from sqlalchemy import text
from sqlalchemy.orm import Session
from models.db import get_db

logger = logging.getLogger(__name__)

api_router = APIRouter(prefix="/api", tags=["user"])


@api_router.get("/health")
async def health_check():
    return {
        "status": "ok",
        "timestamp": datetime.now().isoformat(),
    }


@api_router.post("/food/reserve")
async def reserve_food(request: Request, db: Session = Depends(get_db)):
    body = await request.json()
    order_id = body.get("order_id")

    # Sqlite does not support for update
    # sql_query = text(
    #     "SELECT * FROM packets WHERE is_reserved = false and order_id is null and food_id = :food_id Limit 1 FOR UPDATE SKIP LOCKED"
    # )

    sql_query = text(
        "SELECT * FROM packets WHERE is_reserved = false and order_id is null and food_id = :food_id Limit 1"
    )
    result = db.execute(sql_query, {"food_id": 1})
    rows = result.mappings().all()
    packet_row = {}
    try:
        packet_row = dict(rows[0])
        packet_id = packet_row["id"]
        sql_query = text(
            "UPDATE packets SET is_reserved=true, order_id = :order_id WHERE id = :packet_id"
        )
        result = db.execute(sql_query, {"order_id": order_id, "packet_id": packet_id})
        db.commit()
        return packet_row
    except Exception:
        db.rollback()
        raise HTTPException(status_code=404, detail="Can not reserve food")


@api_router.post("/food/book")
async def book_food(request: Request, db: Session = Depends(get_db)):
    body = await request.json()
    order_id = body.get("order_id")

    sql_query = text(
        "SELECT * FROM packets WHERE is_reserved = true and order_id = :order_id"
    )
    result = db.execute(sql_query, {"order_id": order_id})
    rows = result.mappings().all()
    packet_row = {}
    try:
        packet_row = dict(rows[0])
        packet_id = packet_row["id"]
        sql_query = text("UPDATE packets SET is_reserved=false WHERE id = :packet_id")
        result = db.execute(sql_query, {"packet_id": packet_id})
        db.commit()
        return packet_row
    except Exception:
        db.rollback()
        raise HTTPException(status_code=404, detail="Can not book food")


@api_router.post("/agent/reserve")
async def reserve_agent(request: Request, db: Session = Depends(get_db)):
    body = await request.json()
    order_id = body.get("order_id")

    # sqlite does not support For update skip locked
    # sql_query = text(
    #     "SELECT * FROM delivery_agent WHERE is_reserved = false and order_id is null Limit 1 FOR UPDATE SKIP LOCKED"
    # )

    sql_query = text(
        "SELECT * FROM delivery_agent WHERE is_reserved = false and order_id is null Limit 1"
    )
    result = db.execute(sql_query)
    rows = result.mappings().all()
    agent_row = {}
    try:
        agent_row = dict(rows[0])
        agent_id = agent_row["id"]
        sql_query = text(
            "UPDATE delivery_agent SET is_reserved=true, order_id = :order_id WHERE id = :agent_id"
        )
        result = db.execute(sql_query, {"order_id": order_id, "agent_id": agent_id})
        db.commit()
        return agent_row
    except Exception:
        db.rollback()
        raise HTTPException(status_code=404, detail="Can not reserve agent")


@api_router.post("/agent/book")
async def book_agent(request: Request, db: Session = Depends(get_db)):
    body = await request.json()
    order_id = body.get("order_id")

    sql_query = text(
        "SELECT * FROM delivery_agent WHERE is_reserved = true and order_id = :order_id"
    )
    result = db.execute(sql_query, {"order_id": order_id})
    rows = result.mappings().all()
    packet_row = {}
    try:
        agent_row = dict(rows[0])
        agent_id = agent_row["id"]
        sql_query = text(
            "UPDATE delivery_agent SET is_reserved=false WHERE id = :agent_id"
        )
        result = db.execute(sql_query, {"agent_id": agent_id})
        db.commit()
        return packet_row
    except Exception as e:
        print("Error", str(e))
        db.rollback()
        raise HTTPException(status_code=404, detail="Can not book agent")
