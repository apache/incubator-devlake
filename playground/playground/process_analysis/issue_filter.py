from dataclasses import dataclass
from typing import Optional

import pandas as pd

@dataclass
class IssueFilter:
    project_key: Optional[str] = None
    from_date: Optional[pd.Timestamp] = None
    to_date: Optional[pd.Timestamp] = None
    issue_type: Optional[str] = None

    def apply(self, issue_df: pd.DataFrame ):
        df = issue_df.copy()
        if self.project_key:
            df = df[df['issue_key'].str.startswith(self.project_key)]
        if self.from_date:
            df = df[df['changed_date'] >= self.from_date]
        if self.to_date:
            df = df[df['changed_date'] <= self.to_date]
        if self.issue_type:
            df = df[df['issue_type'] == self.issue_type]
        return df
