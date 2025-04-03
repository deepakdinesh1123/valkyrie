import base64
import json

from contextlib import contextmanager

import websocket


@contextmanager
def websocket_connection(url: str, timeout: int = 30):
    """Context manager for handling WebSocket connections"""
    ws = websocket.WebSocket()
    try:
        ws.connect(url, timeout=timeout)
        yield ws
    finally:
        ws.close()


def decode_websocket_message(encoded_message: str) -> dict:
    # Decode the Base64-encoded message
    decoded_bytes = base64.b64decode(encoded_message)
    decoded_str = decoded_bytes.decode("utf-8")

    # Parse the JSON string
    message_dict = json.loads(decoded_str)

    return message_dict


def encode_websocket_message(message_dict: dict) -> str:
    # Convert the dictionary to a JSON string
    json_str = json.dumps(message_dict)

    # Encode the JSON string in Base64
    encoded_bytes = base64.b64encode(json_str.encode("utf-8"))
    encoded_message = encoded_bytes.decode("utf-8")

    return encoded_message
