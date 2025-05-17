import json
import time

import requests
from dotenv import load_dotenv

from .fixtures import *

load_dotenv()


def test_normal_message(environment, get_mongo):
    print(get_mongo.client.server_info())

    msg_text = {
        "user_id": 1,
        "text": "Hello world!"
    }

    response = requests.post(f'{environment.RECEIVER_URL}/api/message', json=msg_text)
    assert response.status_code == 201

    time.sleep(15)

    output = get_mongo.find_one({"message": msg_text["text"]})
    assert output is not None


def test_spam_message(environment, get_minio):
    msg_text = {
        "user_id": 1,
        "text": "AAAAAAAAAAAAAAAAAAAAAAAAAA"
    }

    response = requests.post(f'{environment.RECEIVER_URL}/api/message', json=msg_text)
    assert response.status_code == 201

    time.sleep(5)

    minio_msg = json.loads(get_minio.get_object())
    assert msg_text['text'] == minio_msg['text']
