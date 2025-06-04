#!/bin/bash

set -e

wait_for_mongo() {
  local host=$1
  echo "⏳ Waiting for $host..."
  until mongosh --host "$host" --quiet --eval "db.adminCommand({ ping: 1 })" > /dev/null 2>&1; do
    sleep 3
  done
  echo "✅ $host is ready."
}

# Ждём все нужные узлы
wait_for_mongo mongo-config:27019
wait_for_mongo mongo-shard1-primary:27018
wait_for_mongo mongo-shard1-secondary:27028
wait_for_mongo mongos:27017

echo "⚙️ Initiating config server replica set..."
mongosh --host mongo-config:27019 --quiet <<EOF
try {
  rs.initiate({
    _id: "configReplSet",
    configsvr: true,
    members: [{ _id: 0, host: "mongo-config:27019" }]
  });
} catch (e) { print(e); }
EOF

echo "⚙️ Initiating shard1 replica set..."
mongosh --host mongo-shard1-primary:27018 --quiet <<EOF
try {
  rs.initiate({
    _id: "shard1",
    members: [
      { _id: 0, host: "mongo-shard1-primary:27018" },
      { _id: 1, host: "mongo-shard1-secondary:27028" }
    ]
  });
} catch (e) { print(e); }
EOF

echo "🔗 Adding shard1 to cluster..."
mongosh --host mongos:27017 --quiet <<EOF
try {
  sh.addShard("shard1/mongo-shard1-primary:27018,mongo-shard1-secondary:27028");
} catch (e) { print(e); }
EOF

echo "🧩 Enabling sharding for database and collection..."
mongosh --host mongos:27017 --quiet <<EOF
try {
  sh.enableSharding("your_db");
  sh.shardCollection("your_db.your_collection", { "user_id": "hashed" });
} catch (e) { print(e); }
EOF

echo "✅ Sharding setup complete."
