FROM openjdk:17-oracle AS builder

COPY . .

RUN ./mvnw clean package
RUN mv target/*.jar application.jar
RUN java -Djarmode=layertools -jar application.jar extract

FROM openjdk:17-oracle

ENV KAFKA_BROKERS localhost:9092
ENV APPLICATION_ID account-bookings-kafka-streams
ENV SCHEMA_REGISTRY_URL http://localhost:8081
ENV KAFKA_TOPIC_ACCOUNTS event_core_banking_accounts
ENV KAFKA_TOPIC_BOOKINGS event_core_banking_bookings
ENV KAFKA_TOPIC_ACCOUNT_BOOKINGS event_core_banking_account_bookings
ENV KAFKA_TOPIC_ACCOUNT_BOOKINGS_PARTITIONS 1
ENV KAFKA_TOPIC_ACCOUNT_BOOKINGS_REPLICAS 1

COPY --from=builder dependencies/ ./
COPY --from=builder snapshot-dependencies/ ./
COPY --from=builder spring-boot-loader/ ./
COPY --from=builder application/ ./

ENTRYPOINT ["java", "org.springframework.boot.loader.JarLauncher"]