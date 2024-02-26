from plotly.graph_objs import Box, Figure

from playground.process_analysis.status_transition_graph import StatusTransitionGraph


class StatusTransitionGraphStatsPlot:
    """Plot statistics of a status transition graph."""

    @staticmethod
    def boxplot(source: StatusTransitionGraph, max_edges: int = 8) -> Figure:
        """Create a boxplot of status transition durations."""
        fig = Figure()

        edges = sorted(source.graph.edges.data(),
                       key=lambda edge: len(edge[2]["durations"]),
                       reverse=True)

        for edge in edges[:max_edges]:
            durations = edge[2]["durations"]
            if len(durations) < 4:
                continue
            fig.add_trace(Box(
                y=durations,
                name=f"{edge[0]} -> {edge[1]}",
                boxpoints='outliers'
            ))

        fig.update_layout(
            title_text="Status Transition Duration Statistics",
            yaxis_title="Days",
            showlegend=False)

        return fig
    