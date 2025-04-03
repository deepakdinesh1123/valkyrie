import logging
import sys


def setup_logging() -> logging.Logger:
    """
    Sets up the logging configuration and returns a custom logger.

    Parameters:
    log_level (str): The log level to filter logs. Should be one of DEBUG, INFO, WARNING, ERROR, CRITICAL.

    Returns:
    logging.Logger: Configured logger instance.
    """
    from .config import config

    numeric_level = getattr(logging, config.LOG_LEVEL.upper(), None)
    if not isinstance(numeric_level, int):
        raise ValueError(f"Invalid log level: {config.LOG_LEVEL}")

    logger = logging.getLogger("pyvalkyrie")
    logger.setLevel(numeric_level)

    # Create handlers
    stream_handler = logging.StreamHandler(sys.stdout)
    stream_handler.setLevel(numeric_level)

    # Create formatters and add it to handlers
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )
    stream_handler.setFormatter(formatter)

    # Add handlers to the logger
    logger.addHandler(stream_handler)

    return logger


logger = setup_logging()
