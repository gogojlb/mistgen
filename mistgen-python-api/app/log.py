import logging
import sys

from app import config

logging.basicConfig(level=config.env.str('LOG_LEVEL'))
LOG = logging.getLogger("API")
LOG.propagate = False

FORMAT = "[%(asctime)s] [%(process)d] [%(levelname)s] %(message)s"
TIMESTAMP_FORMAT = "%Y-%m-%d %H:%M:%S %z"
formatter = logging.Formatter(FORMAT, TIMESTAMP_FORMAT)

stream_handler = logging.StreamHandler(sys.stdout)
stream_handler.setFormatter(formatter)

LOG.addHandler(stream_handler)

def get_logger():
    return LOG