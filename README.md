# Customer Single View

## Project Planning

Gantt Diagram

![Gantt Diagram](planning/gantt.png)


## System Design

![System Design](doc/system-design.png)

## Environment Setup

2 Options:

* [Docker Compose](docker-compose/README.md)
* [K8S-Terraform](k8s-terraform/README.md)

## Data Definition

Sources:

Customer MS (Kafka Events - AVRO)

![Customer MS](doc/data-definition-customer-ms.png)

See: [source](schemas/customer-value.avsc)

Contact Center (Mongo DB - BSON)

![Contact Center](doc/data-definition-contact-center.png)

See: [source](schemas/contact-center.json)

Core Banking (Postgresql - DDL)

![Core Banking](doc/data-definition-core-banking.png)


See: [source](schemas/core-banking.ddl)