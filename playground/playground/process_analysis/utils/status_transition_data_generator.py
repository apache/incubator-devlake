from dataclasses import dataclass
from datetime import datetime, timedelta
from enum import Enum

import numpy.random as rand
import pandas as pd


class Status(Enum):
    BACKLOG = ("Backlog", "TODO")
    TODO = ("To Do", "TODO")
    READY = ("Ready", "TODO")
    IN_PROGRESS = ("In Progress", "IN_PROGRESS")
    REVIEW = ("Review", "IN_PROGRESS")
    TEST = ("Test", "IN_PROGRESS")
    RELEASE = ("Release", "DONE")
    DONE = ("Done", "DONE")
    WONT_FIX = ("Won't Fix", "DONE")


@dataclass(frozen=True)
class StatusChange:
    issue_key: str
    issue_type: str
    created_date: datetime
    from_status: Status
    to_status: Status
    changed_date: datetime


def generate_random_status_changes(n: int = 1000) -> pd.DataFrame:
    changes: list[StatusChange] = []
    count = 0
    status_change = None
    while count < n:
        if status_change is None:
            status_change = _create_first_status_change()
        else:
            status_change = _create_next_status_change(status_change)

        if status_change is not None:
            changes.append(status_change)
            count += 1

    df = pd.DataFrame([{
        "issue_key": change.issue_key,
        "issue_type": change.issue_type,
        "created_date": pd.to_datetime(change.created_date),
        "original_from_value": change.from_status.value[0],
        "from_status": change.from_status.value[1],
        "original_to_value": change.to_status.value[0],
        "to_status": change.to_status.value[1],
        "changed_date": pd.to_datetime(change.changed_date)
    } for change in changes])

    return df


def _create_first_status_change():
    project_key = rand.choice(["CORE", "PLAT", "BRAN", "SUPP", "MISC"], p=[0.4, 0.2, 0.2, 0.1, 0.1])
    issue_key = f"{project_key}-{rand.randint(1, 1000)}"
    issue_type = rand.choice(["Bug", "Task", "Story"], p=[0.3, 0.3, 0.4])
    created_date = datetime(2021, rand.randint(1, 12), rand.randint(1, 28))

    next_status = rand.choice(
                [Status.TODO, Status.READY, Status.IN_PROGRESS, Status.WONT_FIX],
                p=[0.7, 0.15, 0.05, 0.1])
    changed_date = created_date + timedelta(hours=rand.gumbel(24*50, 24*20))
    return StatusChange(issue_key, issue_type, created_date, Status.BACKLOG, next_status, changed_date)


def _create_next_status_change(current: StatusChange) -> StatusChange | None:
    next_status = None
    changed_date = current.changed_date

    match current.to_status:
        case Status.BACKLOG:
            next_status = rand.choice(
                [None, Status.TODO, Status.READY, Status.IN_PROGRESS, Status.WONT_FIX],
                p=[0.1, 0.6, 0.15, 0.05, 0.1])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*50, 24*20))
        case Status.TODO:
            next_status = rand.choice(
                [None, Status.BACKLOG, Status.READY, Status.IN_PROGRESS, Status.WONT_FIX],
                p=[0.1, 0.1, 0.45, 0.25, 0.1])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*28, 24*7))
        case Status.READY:
            next_status = rand.choice(
                [None, Status.BACKLOG, Status.TODO, Status.IN_PROGRESS, Status.WONT_FIX],
                p=[0.1, 0.05, 0.05, 0.75, 0.05])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*14, 24*3))
        case Status.IN_PROGRESS:
            next_status = rand.choice(
                [None, Status.TODO, Status.READY, Status.REVIEW, Status.TEST],
                p=[0.05, 0.05, 0.1, 0.6, 0.2])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*7, 24*3))
        case Status.REVIEW:
            next_status = rand.choice(
                [None, Status.TODO, Status.IN_PROGRESS, Status.TEST],
                p=[0.02, 0.03, 0.15, 0.8])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*1, 24*0.5))
        case Status.TEST:
            next_status = rand.choice(
                [None, Status.TODO, Status.IN_PROGRESS, Status.RELEASE],
                p=[0.02, 0.06, 0.17, 0.75])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*2, 24*1))
        case Status.RELEASE:
            next_status = rand.choice(
                [None, Status.TODO, Status.TEST, Status.DONE],
                p=[0.02, 0.01, 0.1, 0.87])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*3, 24*1))
        case Status.DONE:
            next_status = rand.choice(
                [None, Status.TODO, Status.IN_PROGRESS],
                p=[0.85, 0.02, 0.13])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*7, 24*3))
        case Status.WONT_FIX:
            next_status = rand.choice(
                [None, Status.BACKLOG, Status.TODO, Status.IN_PROGRESS, Status.DONE],
                p=[0.85, 0.05, 0.05, 0.02, 0.03])
            changed_date = current.changed_date + timedelta(hours=rand.gumbel(24*60, 24*10))
        case _:
            pass

    if next_status is None:
        return None

    return StatusChange(current.issue_key, current.issue_type, current.created_date, current.to_status, next_status, changed_date)
