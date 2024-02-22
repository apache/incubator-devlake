from datetime import datetime

import pandas as pd
from playground.process_analysis.status_transition_graph import StatusTransitionGraph
from playground.process_analysis.status_transition_graph_vistualizer import StatusTransitionGraphVisualizer


def test_create_status_transition_graph_from_data_frame():
    result = StatusTransitionGraph.from_data_frame(test_data_frame)

    assert result.total_transition_count == 11
    assert sorted(list(result.graph.nodes.data())) == [
        ('Done', {'count': 2, 'category': 'DONE'}), 
        ('In Progress', {'count': 3, 'category': 'IN_PROGRESS'}), 
        ('Ready', {'count': 3, 'category': 'TODO'}), 
        ('To Do', {'count': 5, 'category': 'TODO'}), 
        ("Won't Fix", {'count': 2, 'category': 'DONE'})
    ]
    assert sorted(list(result.graph.edges.data())) == [
        ('In Progress', 'Done', {'count': 2, 'avg_duration': 1.0}),
        ('In Progress', 'To Do', {'count': 1, 'avg_duration': 4.0}),
        ('Ready', 'In Progress', {'count': 1, 'avg_duration': 2.0}),
        ('To Do', 'In Progress', {'count': 1, 'avg_duration': 1.0}),
        ('To Do', 'Ready', {'count': 3, 'avg_duration': 2.6666666666666665}),
        ('To Do', "Won't Fix", {'count': 2, 'avg_duration': 2.5}),
        ("Won't Fix", 'In Progress', {'count': 1, 'avg_duration': 9.0})
    ]


def test_convert_status_transition_graph_to_graphiz_dot():
    source = StatusTransitionGraph.from_data_frame(test_data_frame)
    dot = StatusTransitionGraphVisualizer().visualize(source)
    
    assert dot.source.replace("\t", "") == """strict digraph "Status Transitions" {
        graph [rankdir=TD]
        subgraph TODO {
            label=TODO
            rank=min
            "To Do" [label=<To Do<BR/>(<FONT POINT-SIZE="8">5x</FONT>)> color=darkslategray fillcolor=lightgray fontname=Arial fontsize=12 penwidth=4.55 style=filled]
            Ready [label=<Ready<BR/>(<FONT POINT-SIZE="8">3x</FONT>)> color=darkslategray fillcolor=lightgray fontname=Arial fontsize=12 penwidth=2.73 style=filled]
        }
        subgraph IN_PROGRESS {
            label=IN_PROGRESS
            rank=same
            "In Progress" [label=<In Progress<BR/>(<FONT POINT-SIZE="8">3x</FONT>)> color=darkslategray fillcolor=yellow fontname=Arial fontsize=12 penwidth=2.73 style=filled]
        }
        subgraph DONE {
            label=DONE
            rank=max
            Done [label=<Done<BR/>(<FONT POINT-SIZE="8">2x</FONT>)> color=darkslategray fillcolor=green fontname=Arial fontsize=12 penwidth=1.82 style=filled]
            "Won\'t Fix" [label=<Won\'t Fix<BR/>(<FONT POINT-SIZE="8">2x</FONT>)> color=darkslategray fillcolor=green fontname=Arial fontsize=12 penwidth=1.82 style=filled]
        }
        "To Do" -> Ready [label=<2.7 days avg<BR/>(<FONT POINT-SIZE="8">3x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=10.91 style=filled]
        "To Do" -> "In Progress" [label=<1.0 days avg<BR/>(<FONT POINT-SIZE="8">1x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=3.64 style=filled]
        "To Do" -> "Won\'t Fix" [label=<2.5 days avg<BR/>(<FONT POINT-SIZE="8">2x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=7.27 style=filled]
        Ready -> "In Progress" [label=<2.0 days avg<BR/>(<FONT POINT-SIZE="8">1x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=3.64 style=filled]
        "In Progress" -> Done [label=<1.0 days avg<BR/>(<FONT POINT-SIZE="8">2x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=7.27 style=filled]
        "In Progress" -> "To Do" [label=<4.0 days avg<BR/>(<FONT POINT-SIZE="8">1x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=3.64 style=filled]
        "Won\'t Fix" -> "In Progress" [label=<9.0 days avg<BR/>(<FONT POINT-SIZE="8">1x</FONT>)> color=darkslategray fontname=Arial fontsize=12 penwidth=3.64 style=filled]
    }
    """.replace("    ", "") # remove indentation


def _pd_timestamp_from(datetime_str: str) -> pd.Timestamp:
    return pd.Timestamp(datetime.strptime(datetime_str, "%Y-%m-%d %H:%M:%S"))

test_data_frame = df = pd.DataFrame([
    {
        "issue_key": "ISSUE-3",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "Ready",
        "to_value": "TODO",
        "changed_date": _pd_timestamp_from("2021-01-07 00:00:00")
    },
    {
        "issue_key": "ISSUE-1",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "Ready",
        "to_value": "TODO",
        "changed_date": _pd_timestamp_from("2021-01-02 00:00:00")
    },
    {
        "issue_key": "ISSUE-1",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "Ready",
        "from_value": "TODO",
        "original_to_value": "In Progress",
        "to_value": "IN_PROGRESS",
        "changed_date": _pd_timestamp_from("2021-01-04 00:00:00")
    },
    {
        "issue_key": "ISSUE-1",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "In Progress",
        "from_value": "IN_PROGRESS",
        "original_to_value": "Done",
        "to_value": "DONE",
        "changed_date": _pd_timestamp_from("2021-01-05 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "In Progress",
        "from_value": "IN_PROGRESS",
        "original_to_value": "Done",
        "to_value": "DONE",
        "changed_date": _pd_timestamp_from("2021-01-18 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "Ready",
        "to_value": "TODO",
        "changed_date": _pd_timestamp_from("2021-01-02 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "In Progress",
        "to_value": "IN_PROGRESS",
        "changed_date": _pd_timestamp_from("2021-01-03 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "In Progress",
        "from_value": "IN_PROGRESS",
        "original_to_value": "To Do",
        "to_value": "TODO",
        "changed_date": _pd_timestamp_from("2021-01-07 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "Won't Fix",
        "to_value": "DONE",
        "changed_date": _pd_timestamp_from("2021-01-08 00:00:00")
    },
    {
        "issue_key": "ISSUE-2",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "Won't Fix",
        "from_value": "DONE",
        "original_to_value": "In Progress",
        "to_value": "IN_PROGRESS",
        "changed_date": _pd_timestamp_from("2021-01-17 00:00:00")
    },
    {
        "issue_key": "ISSUE-4",
        "created_date": _pd_timestamp_from("2021-01-01 00:00:00"),
        "original_from_value": "To Do",
        "from_value": "TODO",
        "original_to_value": "Won't Fix",
        "to_value": "DONE",
        "changed_date": _pd_timestamp_from("2021-01-05 00:00:00")
    },
])
