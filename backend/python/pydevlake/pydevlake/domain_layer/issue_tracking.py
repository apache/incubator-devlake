from datetime import datetime
from enum import Enum
from sqlmodel import Field

from pydevlake.model import DomainModel, DomainScope, NoPKModel


class Board(DomainScope, table=True):
    __tablename__ = "boards"

    created_at: datetime
    updated_at: datetime
    _raw_data_params: str
    _raw_data_table: str
    _raw_data_id: int
    _raw_data_remark: str
    name: str
    description: str
    url: str
    created_date: datetime
    type: str


class BoardIssue(NoPKModel, table=True):
    __tablename__ = "board_issues"

    board_id: str = Field(primary_key=True)
    issue_id: str = Field(primary_key=True)


class Issue(DomainModel, table=True):
    class Type(Enum):
        Incident = 'INCIDENT'

    class Status(Enum):
        Done = 'DONE'
        InProgress = 'IN_PROGRESS'
        ToDo = 'TODO'

    __tablename__ = "issues"

    url: str
    icon_url: str
    issue_key: str
    title: str
    description: str
    epic_key: str
    type: str
    original_type: str
    status: str
    original_status: str
    resolution_date: datetime
    created_date: datetime
    updated_date: datetime
    lead_time_minutes: int
    parent_issue_id: str
    priority: str
    story_point: float
    original_estimate_minutes: int
    time_spent_minutes: int
    time_remaining_minutes: int
    creator_id: str
    creator_name: str
    assignee_id: str
    assignee_name: str
    severity: str
    component: str
    original_project: str
    urgency: str
