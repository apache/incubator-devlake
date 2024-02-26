from dataclasses import dataclass, field
from itertools import groupby
from operator import itemgetter
import enum

import graphviz
import statistics

from playground.process_analysis.status_transition_graph import StatusTransitionGraph

class StatisticLabelConfig(enum.Enum):
    AVG = "avg"
    MEDIAN = "med"
    IQR = "iqr"
    MIN_MAX = "min-max"
    

@dataclass
class VisualizerConfig:
    font: str = "Arial"
    fontsize: str = "12"
    sub_fontsize: str = "8"
    node_border_color: str = "darkslategray"
    edge_color: str = "darkslategray"
    category_fill_color: dict[str, str] = field(
        default_factory=lambda: {
            "TODO": "lightgray",
            "IN_PROGRESS": "yellow",
            "DONE": "green",
        }
    )
    fallback_fill_color: str = "aliceblue"
    node_penwidth_factor: float = 10.0
    edge_penwidth_factor: float = 40.0


class StatusTransitionGraphVisualizer:
    """Visualize a status transition graph using Graphviz."""

    def __init__(self, config: VisualizerConfig | None = None) -> None:
        if config is None:
            config = VisualizerConfig()
        self.config = config

    def visualize(
        self, 
        source: StatusTransitionGraph, 
        threshold: float = 1.0,
        label_statistic: StatisticLabelConfig = StatisticLabelConfig.AVG
    ) -> graphviz.Digraph:
        """Create a Graphviz digraph from a StatusTransitionGraph.
        
        Args:
            source: The StatusTransitionGraph to visualize.
            threshold: Number between 0.0 and 1.0, defaults to 1.0 (or 100%.)
                Exclude edges from the visualization that represent less than 
                the percentage of total status transition within the threshold.
        """

        dot_graph = graphviz.Digraph("Status Transitions", format="svg", strict=True)
        dot_graph.attr("graph", rankdir="TD")
        dot_graph.attr("node", color=self.config.node_border_color, **self.__default_attrs())
        dot_graph.attr("edge", color=self.config.edge_color, **self.__default_attrs())

        graph = source.graph
        for category, nodes in groupby(graph.nodes(data="category"), itemgetter(1)):
            with dot_graph.subgraph(name=f"{category}") as cluster:
                cluster.attr(label=str(category))
                cluster.attr(rank=self.__category_rank(category))
                cluster.attr("node", fillcolor=self.__category_color(category))
                for node, _ in nodes:
                    count = graph.nodes[node]["count"]
                    penwidth = count / source.total_transition_count * self.config.node_penwidth_factor
                    cluster.node(
                        name=node,
                        label=self.__node_label(node, count),
                        penwidth=str(round(penwidth, 2)),
                    )

        for edge in graph.edges.data():
            to_count = len(edge[2]["durations"])
            if to_count > (1.00 - threshold) * source.total_transition_count:
                penwidth = to_count / source.total_transition_count * self.config.edge_penwidth_factor
                dot_graph.edge(
                    edge[0],
                    edge[1],
                    labeltooltip=self.__edge_tooltip(edge[2]["durations"]),
                    label=self.__edge_label(edge[2]["durations"], label_statistic),
                    penwidth=str(round(penwidth, 2)),
                )

        return dot_graph

    @staticmethod
    def is_dot_executable_available() -> bool:
        """Check if the 'dot' executable is available."""
        try:
            graphviz.version()
            return True
        except graphviz.ExecutableNotFound:
            return False

    def __category_rank(self, category: str) -> str:
        match category:
            case "TODO":
                return "min"
            case "DONE":
                return "max"
            case _:
                return ""

    def __category_color(self, category: str) -> str:
        return self.config.category_fill_color.get(
            category, self.config.fallback_fill_color
        )

    def __node_label(self, name: str, count: int) -> str:
        return f'<{name}<BR/><FONT POINT-SIZE="{self.config.sub_fontsize}">({str(count)+"x"})</FONT>>'

    def __edge_label(self, durations: list[float], label_statistic: StatisticLabelConfig) -> str:
        stat = ''
        match label_statistic:
            case StatisticLabelConfig.AVG:
                stat = f'{statistics.mean(durations):0.1f} days avg'
            case StatisticLabelConfig.MEDIAN:
                stat = f'{statistics.median(durations):0.1f} days med'
            case StatisticLabelConfig.IQR:
                if len(durations) >= 4:
                    quantiles = statistics.quantiles(durations, n=4)
                    stat = f'{quantiles[0]:0.1f} - {quantiles[2]:0.1f} days IQR (25-75%)'
            case StatisticLabelConfig.MIN_MAX:
                if len(durations) >= 2:
                    stat = f'{min(durations):0.1f} - {max(durations):0.1f} days min-max'
        count_str = f'<FONT POINT-SIZE="{self.config.sub_fontsize}">({str(len(durations))+"x"})</FONT>'
        return f'<{stat}<BR/>{count_str}>'

    def __edge_tooltip(self, durations: list[float]) -> str:
        lines = []
        lines.append(f"avg: {statistics.mean(durations):0.1f} days")
        lines.append(f"med: {statistics.median(durations):0.1f} days")
        if len(durations) >= 4:
            quantiles = statistics.quantiles(durations, n=4)
            lines.append(f"IQR (25-75%): {quantiles[0]:0.1f} - {quantiles[2]:0.1f} days")
        if len(durations) >= 2:
            lines.append(f"min-max: {min(durations):0.1f} - {max(durations):0.1f} days")
        return "\n".join(lines)

    def __default_attrs(self) -> dict:
        return {
            "style": "filled",
            "fontname": self.config.font,
            "fontsize": self.config.fontsize,
        }
