docker compose down -v

Stops and removes:

All running containers

All Docker networks created by compose

All volumes (-v) → Database data will be deleted

👉 Use this when:

Database is corrupted

Kafka/Zookeeper not starting

Port conflicts

You want a fresh environment

⚠️ Warning: This deletes PostgreSQL data.