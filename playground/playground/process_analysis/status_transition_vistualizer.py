from itertools import groupby
from operator import itemgetter

import graphviz

from playground.process_analysis.status_transition_graph import StatusTransitionGraph

FONT = "Arial"
FONTSIZE = "12"
FONTSIZE_S = "8"
NODE_BORDER_COLOR = "darkgray"
EDGE_COLOR = "darkslategray"
FILL_COLOR_MAP = {"TODO": "lightgray", "IN_PROGRESS": "yellow", "DONE": "green"}
FALLBACK_COLOR = "aliceblue"

class StatusTransitionVisualizer:
    @staticmethod
    def visualize(source: StatusTransitionGraph, threshold: float = 1.0) -> graphviz.Digraph:
        dot_graph = graphviz.Digraph("Status Transitions", format="svg", strict=True)
        dot_graph.attr("graph", rankdir="TD")

        graph = source.graph
        for category, nodes in groupby(
            graph.nodes(data="category"), itemgetter(1)
        ):
            with dot_graph.subgraph(name=f"sub_{category}") as cluster:
                cluster.attr(label=str(category))
                cluster.attr(rank=category_rank(category))
                for item in nodes:
                    node_name, _ = item
                    node = graph.nodes[node_name]
                    cluster.node(
                        name=node_name,
                        label=node_label(node_name, node["count"]),
                        color=NODE_BORDER_COLOR,
                        fillcolor=category_color(category),
                        penwidth=str(node["count"] / source.total_transition_count * 10),
                        **default_attrs(),
                    )

        for edge in graph.edges(data=True):
            if edge[2]["count"] > (1.00 - threshold) * source.total_transition_count:
                dot_graph.edge(
                    edge[0],
                    edge[1],
                    label=edge_label(edge[2]["avg_duration"], edge[2]["count"]),
                    penwidth=str(edge[2]["count"] / source.total_transition_count * 40),
                    color=EDGE_COLOR,
                    **default_attrs(),
                )

        return dot_graph


def category_rank(category: str) -> str:
    match category:
        case "TODO":
            return "min"
        case "DONE":
            return "max"
        case _:
            return "same"


def category_color(category: str) -> str:
    return FILL_COLOR_MAP.get(category, FALLBACK_COLOR)


def node_label(name: str, count: int) -> str:
    return f'<{name}<BR/>(<FONT POINT-SIZE="{FONTSIZE_S}">{str(count)+"x"}</FONT>)>'


def edge_label(avg_duration: float, count: int) -> str:
    return f'<{avg_duration:0.1f} days avg<BR/>(<FONT POINT-SIZE="{FONTSIZE_S}">{str(count)+"x"}</FONT>)>'


def default_attrs() -> dict:
    return {
        "style": "filled",
        "fontname": FONT,
        "fontsize": FONTSIZE,
    }
