# Kafka Docker Setup for E2E Testing

## Table of Contents
- [Redpanda (Recommended)](#redpanda-recommended)
- [Apache Kafka KRaft](#apache-kafka-kraft)
- [Comparison](#comparison)
- [Health Checks](#health-checks)
- [Testcontainers Alternative](#testcontainers-alternative)

---

## Redpanda (Recommended)

Redpanda is a Kafka-compatible streaming platform that starts 2x faster and requires no ZooKeeper.

```yaml
# docker-compose.e2e.yml
version: '3.8'
services:
  redpanda:
    image: redpandadata/redpanda:latest
    hostname: redpanda
    ports:
      - 9092:9092   # Internal broker port
      - 9094:9094   # External broker port (for test clients)
      - 18082:18082 # Pandaproxy (REST API)
    command:
      - redpanda start
      - --kafka-addr internal://0.0.0.0:9092,external://0.0.0.0:9094
      - --advertise-kafka-addr internal://redpanda:9092,external://localhost:9094
      - --mode dev-container
      - --smp 1
      - --memory 512M
      - --overprovisioned
    healthcheck:
      test: ["CMD-SHELL", "rpk cluster health | grep -q 'Healthy' || exit 1"]
      interval: 5s
      timeout: 10s
      retries: 10
      start_period: 10s

  redpanda-console:
    image: redpandadata/console:latest
    ports:
      - 8082:8080
    environment:
      KAFKA_BROKERS: redpanda:9092
    depends_on:
      - redpanda

volumes:
  redpanda-data:
    driver: local
```

### Why Redpanda?
- **Faster startup**: ~5.1s vs ~8.5s for Kafka
- **No ZooKeeper**: Single binary, simpler setup
- **Lower memory**: C++ implementation vs JVM
- **Kafka-compatible**: Drop-in replacement API
- **Built-in UI**: Console included

---

## Apache Kafka KRaft

Use when you need exact Kafka compatibility or specific Kafka features.

```yaml
# docker-compose.e2e.yml
version: '3.8'

services:
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    hostname: kafka
    container_name: kafka-e2e
    ports:
      - "9092:9092"
      - "9094:9094"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: 'CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT'
      KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9094'
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
      KAFKA_PROCESS_ROLES: 'broker,controller'
      KAFKA_CONTROLLER_QUORUM_VOTERS: '1@kafka:29093'
      KAFKA_LISTENERS: 'PLAINTEXT://kafka:29092,CONTROLLER://kafka:29093,PLAINTEXT_HOST://0.0.0.0:9094'
      KAFKA_INTER_BROKER_LISTENER_NAME: 'PLAINTEXT'
      KAFKA_CONTROLLER_LISTENER_NAMES: 'CONTROLLER'
      KAFKA_LOG_DIRS: '/tmp/kraft-combined-logs'
      CLUSTER_ID: 'MkU3OEVBNTcwNTJENDM2Qk'
      # Performance tuning for tests
      KAFKA_LOG_RETENTION_MS: 300000        # 5 minutes
      KAFKA_LOG_SEGMENT_BYTES: 1048576      # 1MB segments
    healthcheck:
      test: ["CMD-SHELL", "kafka-broker-api-versions --bootstrap-server localhost:9092 || exit 1"]
      interval: 5s
      timeout: 10s
      retries: 10
      start_period: 30s

volumes:
  kafka-e2e-data:
    driver: local
```

### Apache Kafka Native Image (Faster)

For 40% faster startup, use the native image:

```yaml
services:
  kafka:
    image: apache/kafka-native:latest
    # ... same config as above
```

---

## Comparison

| Metric | Apache Kafka | Kafka Native | Redpanda |
|--------|--------------|--------------|----------|
| Container startup | ~8.5s | ~5.1s | ~5.1s |
| Memory usage | High (JVM) | Medium | Low (C++) |
| ZooKeeper required | No (KRaft) | No (KRaft) | No |
| Kafka compatibility | 100% | 100% | ~99% |
| Best for | Production parity | Faster CI | Fast dev/test |

---

## Health Checks

### Redpanda Health Check Script

```bash
#!/bin/bash
# wait-for-redpanda.sh
TIMEOUT=${1:-30}
HOST=${2:-localhost}
PORT=${3:-9094}

echo "Waiting for Redpanda at $HOST:$PORT..."
for i in $(seq 1 $TIMEOUT); do
  if docker-compose -f docker-compose.e2e.yml exec redpanda \
    rpk cluster health 2>/dev/null | grep -q "Healthy"; then
    echo "Redpanda is ready!"
    exit 0
  fi
  sleep 1
done
echo "Timeout waiting for Redpanda"
exit 1
```

### Kafka Health Check Script

```bash
#!/bin/bash
# wait-for-kafka.sh
TIMEOUT=${1:-60}
HOST=${2:-localhost}
PORT=${3:-9092}

echo "Waiting for Kafka at $HOST:$PORT..."
for i in $(seq 1 $TIMEOUT); do
  if docker-compose -f docker-compose.e2e.yml exec kafka \
    kafka-broker-api-versions --bootstrap-server localhost:9092 2>/dev/null; then
    echo "Kafka is ready!"
    exit 0
  fi
  sleep 1
done
echo "Timeout waiting for Kafka"
exit 1
```

---

## Testcontainers Alternative

For programmatic container management in tests:

```typescript
import { KafkaContainer, StartedKafkaContainer } from '@testcontainers/kafka';

let kafkaContainer: StartedKafkaContainer;

beforeAll(async () => {
  kafkaContainer = await new KafkaContainer('redpandadata/redpanda:latest')
    .withExposedPorts(9092)
    .start();

  const bootstrapServers = kafkaContainer.getBootstrapServers();
  // Use bootstrapServers for KafkaTestHelper
}, 60000);

afterAll(async () => {
  await kafkaContainer?.stop();
});
```

### Container Reuse (Development Only)

```typescript
// For faster local development - reuse containers between runs
const kafka = new KafkaContainer('redpandadata/redpanda:latest')
  .withReuse();

// Requires ~/.testcontainers.properties:
// testcontainers.reuse.enable=true
```

**Warning:** Don't use container reuse in CI - it can cause test pollution.

### Parallel Container Startup

```typescript
// Start multiple containers in parallel for faster setup
beforeAll(async () => {
  const [kafkaContainer, mongoContainer, redisContainer] = await Promise.all([
    new KafkaContainer('redpandadata/redpanda:latest').start(),
    new MongoDBContainer('mongo:6').start(),
    new GenericContainer('redis:7').withExposedPorts(6379).start(),
  ]);

  // Sequential: 15-20s, Parallel: 6-8s
}, 60000);
```

---

## Broker Configuration for Tests

### Recommended Settings

```yaml
environment:
  # Reduce rebalance delay for faster tests
  KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0

  # Faster log cleanup
  KAFKA_LOG_RETENTION_MS: 300000      # 5 minutes
  KAFKA_LOG_SEGMENT_BYTES: 1048576    # 1MB

  # Single replica for speed (test only!)
  KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
  KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
```

### Listener Configuration

For tests connecting from the host machine:

```yaml
# Internal (container-to-container): port 9092
# External (host-to-container): port 9094
KAFKA_LISTENERS: 'PLAINTEXT://kafka:29092,CONTROLLER://kafka:29093,PLAINTEXT_HOST://0.0.0.0:9094'
KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9094'
```

In your tests, use `localhost:9094` as the broker address.
