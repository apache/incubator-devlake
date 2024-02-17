from urllib.parse import urlparse, parse_qsl
from sqlalchemy.engine import Engine, create_engine

def create_db_engine(db_url: str) -> Engine:
    # SQLAlchemy doesn't understand postgres:// scheme
    db_url = db_url.replace("postgres://", "postgresql://")
    # Use MySQL connector for mysql:// scheme
    db_url = db_url.replace("mysql://", "mysql+mysqlconnector://")
    # Remove query args
    base_url = db_url.split('?')[0]
    # `parseTime` parameter is not understood by MySQL driver,
    # so we have to parse query args to remove it
    connect_args = dict(parse_qsl(urlparse(db_url).query))
    if 'parseTime' in connect_args:
        del connect_args['parseTime']
    if 'loc' in connect_args:
        del connect_args['loc']
    if 'tls' in connect_args:
        del connect_args['tls']
        connect_args['ssl'] = {'verify_cert': 'False'}
    try:
        return create_engine(base_url, connect_args=connect_args)
    except Exception as e:
        raise IOError("Unable to make a database connection") from e