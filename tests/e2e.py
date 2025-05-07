import requests
from dotenv import load_dotenv

from .fixtures import *

load_dotenv()


def test_normal_message(environment, get_kafka, get_mongo):
    msg_text = {
        "user_id": 1,
        "text": "Hello world!"
    }

    response = requests.post(f'{environment.RECEIVER_URL}/api/message', json=msg_text)
    assert response.status_code == 201

    kafka_msg = get_kafka.read_message()
    assert msg_text["text"] in kafka_msg

    output = get_mongo.find_one({"text": msg_text["text"]})
    assert output is not None


def test_spam_message(environment, get_kafka, get_minio):
    msg_text = {
        "user_id": 1,
        "text": "AAAAAAAAAAAAAAAAAAAAAAAAAA"
    }

    response = requests.post(f'{environment.RECEIVER_URL}/api/message', json=msg_text)
    assert response.status_code == 201

    kafka_msg = get_kafka.read_message()
    assert msg_text["text"] in kafka_msg

    minio_msg = get_minio.get_object()
    assert msg_text == minio_msg


def test_end_to_end(environment, get_kafka, get_mongo, get_minio):
    test_normal_message(environment, get_kafka, get_mongo)
    test_spam_message(environment, get_kafka, get_minio)