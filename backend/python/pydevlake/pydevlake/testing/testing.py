import pytest

from typing import Union, Type, Iterable

from pydevlake.plugin import Plugin
from pydevlake.model import DomainModel


def assert_convert(plugin: Union[Plugin, Type[Plugin]], stream_name: str, raw: dict, expected: Union[DomainModel, Iterable[DomainModel]]):
    if isinstance(plugin, type):
        plugin = plugin()
    stream = plugin.get_stream(stream_name)
    tool_model = stream.extract(raw)
    domain_models = stream.convert(tool_model, None)
    if not isinstance(expected, Iterable):
        expected = [expected]
    if not isinstance(domain_models, Iterable):
        domain_models = [domain_models]
    for res, exp in zip(domain_models, expected):
        assert res == exp
