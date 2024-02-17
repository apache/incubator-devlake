from datetime import timedelta
from dataclasses import dataclass

import pandas as pd
import networkx as nx

from sqlalchemy.engine import Engine


@dataclass(frozen=True)
class StatusChange:
    issue_key: str
    created_date: pd.Timestamp
    original_from_value: str
    from_value: str
    original_to_value: str
    to_value: str
    changed_date: pd.Timestamp


class StatusTransitionGraph:
    def __init__(self):
        self.graph = nx.DiGraph()
        self.total_transition_count = 0

    def add_status_change(
        self, status_change: StatusChange, previous_status_change: StatusChange | None
    ):
        self.__update_nodes(status_change, previous_status_change)
        self.__update_edges(status_change, previous_status_change)
        self.total_transition_count += 1

    def __update_nodes(
        self, status_change: StatusChange, previous_status_change: StatusChange | None
    ):
        from_status = status_change.original_from_value
        if not self.graph.has_node(from_status):
            self.graph.add_node(from_status, count=0, category=status_change.from_value)

        to_status = status_change.original_to_value
        if not self.graph.has_node(to_status):
            self.graph.add_node(to_status, count=0, category=status_change.to_value)

        if not for_same_issue(status_change, previous_status_change):
            self.graph.nodes[from_status]["count"] += 1
        self.graph.nodes[to_status]["count"] += 1

    def __update_edges(
        self, status_change: StatusChange, previous_status_change: StatusChange | None
    ):
        duration = days_between(status_change, previous_status_change)
        edge_from = status_change.original_from_value
        edge_to = status_change.original_to_value
        if self.graph.has_edge(edge_from, edge_to):
            self.graph.edges[edge_from, edge_to]["avg_duration"] = (
                calculate_avg_duration(
                    self.graph.edges[edge_from, edge_to]["count"],
                    self.graph.edges[edge_from, edge_to]["avg_duration"],
                    duration,
                )
            )
            self.graph.edges[edge_from, edge_to]["count"] += 1
        else:
            self.graph.add_edge(edge_from, edge_to, count=1, avg_duration=duration)

    @classmethod
    def from_database(cls, db_engine: Engine) -> "StatusTransitionGraph":
        query = "select i.issue_key as issue_key, i.created_date as created_date, \
                    ic.original_from_value as original_from_value, ic.from_value as from_value, \
                    ic.original_to_value as original_to_value, ic.to_value as to_value, \
                    ic.created_date as changed_date \
                from issue_changelogs ic \
                    join issues i on i.id = ic.issue_id \
                where ic.field_name = 'status';"
        df = pd.read_sql_query(query, db_engine)
        return cls.from_data_frame(df)

    @classmethod
    def from_data_frame(cls, df: pd.DataFrame) -> "StatusTransitionGraph":
        process_graph: StatusTransitionGraph = cls()

        df = df.copy().sort_values(by=["issue_key", "changed_date"], ascending=True)

        previous_status_change: StatusChange = None
        for item in df.itertuples(index=False):
            status_change = StatusChange(*item)
            process_graph.add_status_change(status_change, previous_status_change)
            previous_status_change = status_change

        return process_graph


def for_same_issue(
    status_change: StatusChange, previous_status_change: StatusChange | None
) -> bool:
    if previous_status_change is None:
        return False
    return status_change.issue_key == previous_status_change.issue_key


def days_between(
    status_change: StatusChange, previous_status_change: StatusChange | None
) -> float:
    if for_same_issue(status_change, previous_status_change):
        return timedelta_between(
            status_change.changed_date, previous_status_change.changed_date
        ) / timedelta(days=1)
    return timedelta_between(
        status_change.changed_date, status_change.created_date
    ) / timedelta(days=1)


def timedelta_between(current: pd.Timestamp, previous: pd.Timestamp) -> timedelta:
    return current.to_pydatetime() - previous.to_pydatetime()


def calculate_avg_duration(count: int, avg_duration: float, duration: float) -> float:
    return (avg_duration * count + duration) / (count + 1)