from sqlalchemy import Column, Integer, String, DateTime, Boolean
from datetime import datetime
from models.db import Base


class Food(Base):
    __tablename__ = "food"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)

    created_at = Column(DateTime, default=datetime.utcnow)


class Packets(Base):
    __tablename__ = "packets"

    id = Column(Integer, primary_key=True, index=True)

    food_id = Column(Integer, nullable=False)
    is_reserved = Column(Boolean, default=False)
    order_id = Column(String, nullable=True)

    created_at = Column(DateTime, default=datetime.utcnow)


class DeliveryAgent(Base):
    __tablename__ = "delivery_agent"

    id = Column(Integer, primary_key=True, index=True)
    is_reserved = Column(Boolean, default=False)
    order_id = Column(String, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)
