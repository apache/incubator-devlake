import pytest
pytest.register_assert_rewrite('pydevlake.testing')

from .testing import assert_convert
