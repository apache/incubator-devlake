from typing import Type

from sqlalchemy import sql
from sqlmodel import Session

from pydevlake import Context


def get(ctx: Context, model_type: Type, *query) -> any:
    with Session(ctx.engine) as session:
        stmt = sql.select(model_type).filter(*query)
        model = session.exec(stmt).scalar()
        return model
