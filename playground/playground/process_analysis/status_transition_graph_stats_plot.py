from plotly.graph_objs import Box, Figure

from playground.process_analysis.status_transition_graph import StatusTransitionGraph


class StatusTransitionGraphStatsPlot:
    """Plot statistics of a status transition graph."""

    @staticmethod
    def boxplot(source: StatusTransitionGraph, max_edges: int = 8) -> Figure:
        """Create a boxplot of status transition durations."""
        fig = Figure()

        most_common_edges_upto_max = sorted(source.graph.edges.data(),
                       key=lambda edge: len(edge[2]["durations"]),
                       reverse=True)[:max_edges]

        edge_counts = list(map(lambda edge: len(edge[2]["durations"]), most_common_edges_upto_max))
        lowest_edge_count = min(edge_counts, default=1)
        highest_edge_count = max(edge_counts, default=1)

        edges_by_max_duration = sorted(most_common_edges_upto_max,
                       key=lambda edge: max(edge[2]["durations"]),
                       reverse=True)

        for edge in edges_by_max_duration:
            durations = list(map(lambda d: round(d, 5), edge[2]["durations"]))
            color = _color_for_edge(len(durations), lowest_edge_count, highest_edge_count)
            fig.add_trace(Box(
                y=durations,
                name=f"{edge[0]} ⮕ {edge[1]} ({len(durations)}x)",
                boxpoints='outliers',
                boxmean=True,
                marker_color=color,
                line_color=color,
            ))

        fig.update_layout(
            title_text="Status Transition Duration Statistics",
            yaxis_title="Days",
            xaxis_title="Transitions with occurences",
            showlegend=False)

        return fig


def _color_for_edge(edge_count: int, min_count: int, max_count: int) -> str:
    base_r, base_g, base_b = (106, 29, 87) # dark base color
    factor = 0 if min_count == max_count else \
        (1 - (edge_count - min_count) / (max_count - min_count)) * 0.8 # 0.8 is the max brightning factor

    edge_r = base_r + (255 - base_r) * factor
    edge_g = base_g + (255 - base_g) * factor
    edge_b = base_b + (255 - base_b) * factor

    return f'rgb({edge_r},{edge_g},{edge_b})'
