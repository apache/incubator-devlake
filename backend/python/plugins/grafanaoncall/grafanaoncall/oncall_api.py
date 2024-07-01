import os

from typing import Optional
from urllib.parse import parse_qs, urlparse

from grafanaoncall.models import GrafanaOncallConnection
from pydevlake.api import API, Paginator, Request, request_hook


class GrafanaOncallPaginator(Paginator):
    def get_items(self, response) -> Optional[list[object]]:
        return response.json['results']

    def get_next_page_id(self, response) -> Optional[str]:
        if not (next_url := response.json.get('next')):
            return

        parsed_url = urlparse(next_url)
        return parse_qs(parsed_url.query)['page'][0]

    def set_next_page_param(self, request, next_page_id):
        request.query_args['page'] = next_page_id


class GrafanaOncallAPI(API):
    paginator = GrafanaOncallPaginator()
    base_url: str = ''
    connection: GrafanaOncallConnection

    def __init__(self, connection: GrafanaOncallConnection):
        super().__init__(connection)
        self.base_url = connection.endpoint

    @request_hook
    def authenticate(self, request: Request):
        request.headers['Authorization'] = self.connection.token.get_secret_value()

    @request_hook
    def verify(self, request: Request):
        if conn_verify := self.connection.verify:
            request.verify = {'true': True, 'false': False}.get(conn_verify, conn_verify)  # true, false or column value
        elif os.getenv('IN_SECURE_SKIP_VERIFY'):
            request.verify = False

    def alert_groups(self):
        return self.get('alert_groups')

    def escalation_chains(self):
        return self.get('escalation_chains')

    def routes(self):
        return self.get('routes')

    def users(self):
        return self.get('users')
