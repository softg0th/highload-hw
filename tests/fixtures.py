import io
import os
from urllib.parse import quote_plus

import pytest

from confluent_kafka import Consumer, KafkaException
from minio import Minio
from pymongo import MongoClient


class Environment:
    def __init__(
            self,
            receiver_url: str,
            bootstrap_server: str,
            minio_endpoint: str,
            access_key_id: str,
            secret_access_key: str,
            minio_bucket: str,
            db_url: str,
            db_name: str,
            collection_name: str,
            mongo_user: str,
            mongo_password: str,
    ):
        self.RECEIVER_URL = receiver_url
        self.BOOTSTRAP_SERVER = bootstrap_server
        self.MINIO_ENDPOINT = minio_endpoint
        self.ACCESS_KEY_ID = access_key_id
        self.SECRET_ACCESS_KEY = secret_access_key
        self.MINIO_BUCKET = minio_bucket
        self.DB_URL = db_url
        self.DB_NAME = db_name
        self.COLLECTION_NAME = collection_name
        self.MONGO_INITDB_ROOT_USERNAME = mongo_user
        self.MONGO_INITDB_ROOT_PASSWORD = mongo_password


def init_env():
    environment = Environment(
        os.getenv('RECEIVER_URL'),
        os.getenv('BOOTSTRAP_SERVER'),
        os.getenv('MINIO_ENDPOINT'),
        os.getenv('ACCESS_KEY_ID'),
        os.getenv('SECRET_ACCESS_KEY'),
        os.getenv('MINIO_BUCKET'),
        os.getenv('DB_URL'),
        os.getenv('DB_NAME'),
        os.getenv('COLLECTION_NAME'),
        os.getenv('MONGO_INITDB_ROOT_USERNAME'),
        os.getenv('MONGO_INITDB_ROOT_PASSWORD')
    )
    return environment


class KafkaTestConsumer:
    def __init__(self, bootstrap_servers: str, topic: str):
        self.topic = topic
        self.consumer = Consumer({
            'bootstrap.servers': bootstrap_servers,
            'group.id': 'test-group',
            'auto.offset.reset': 'earliest'
        })
        self.consumer.subscribe([self.topic])

    def read_message(self, timeout=10.0):
        msg = self.consumer.poll(timeout)
        if msg is None:
            return None
        if msg.error():
            raise KafkaException(msg.error())
        return msg.value().decode('utf-8')

    def close(self):
        self.consumer.close()


class MinioClient:
    def __init__(self, endpoint, access_key, secret_key, bucket):
        self.client = Minio(endpoint, access_key=access_key, secret_key=secret_key, secure=False)
        self.bucket = bucket

    def ensure_bucket(self):
        if not self.client.bucket_exists(self.bucket):
            self.client.make_bucket(self.bucket)

    def put_object(self, object_name: str, data: bytes):
        self.ensure_bucket()
        stream = io.BytesIO(data)
        self.client.put_object(
            self.bucket,
            object_name,
            stream,
            length=len(data),
            content_type="text/plain"
        )

    def get_object(self, object_name: str = None) -> str:
        if object_name is None:
            objects = self.client.list_objects(self.bucket)
            object_name = next(objects).object_name
        response = self.client.get_object(self.bucket, object_name)
        return response.read().decode("utf-8")


class Mongo:
    def __init__(self, db_url, db_name, collection_name):
        user = 'root'
        password = 'example'

        if "@" not in db_url:
            db_url = f"mongodb://{user}:{password}@{db_url.lstrip('mongodb://')}/?authSource=admin"

        self.client = MongoClient(db_url)
        self.collection = self.client[db_name][collection_name]

    def find_one(self, query: dict) -> dict | None:
        return self.collection.find_one(query)

@pytest.fixture(scope="session")
def environment():
    return init_env()

@pytest.fixture(scope="session")
def get_kafka(environment):
    kafka = KafkaTestConsumer(environment.BOOTSTRAP_SERVER, topic="test")
    yield kafka
    kafka.close()

@pytest.fixture(scope="session")
def get_mongo(environment):
    return Mongo(
        environment.DB_URL,
        environment.DB_NAME,
        environment.COLLECTION_NAME
    )

@pytest.fixture(scope="session")
def get_minio(environment):
    return MinioClient(
        environment.MINIO_ENDPOINT,
        environment.ACCESS_KEY_ID,
        environment.SECRET_ACCESS_KEY,
        environment.MINIO_BUCKET
    )
