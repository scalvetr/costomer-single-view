package edu.uoc.scalvetr.accountbookings;

import edu.uoc.scalvetr.Account;
import edu.uoc.scalvetr.AccountBookings;
import edu.uoc.scalvetr.Booking;
import io.confluent.kafka.serializers.AbstractKafkaSchemaSerDeConfig;
import io.confluent.kafka.streams.serdes.avro.SpecificAvroSerde;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.admin.NewTopic;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.kstream.*;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Import;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.annotation.EnableKafkaStreams;
import org.springframework.kafka.config.TopicBuilder;

import java.time.Duration;
import java.util.Map;

@SpringBootApplication
@EnableConfigurationProperties
@EnableKafka
@EnableKafkaStreams
@Import(KafkaStreamsApplicationProperties.class)
@RequiredArgsConstructor
@Slf4j
public class AccountBookingsKafkaStreamsApplication {

    public static void main(String[] args) {
        SpringApplication.run(AccountBookingsKafkaStreamsApplication.class, args);
    }

    private final KafkaStreamsApplicationProperties properties;

    @Bean
    NewTopic output() {
        return TopicBuilder.name(properties.getTopicAccountBookings())
                // .partitions(1).replicas(1) rely on the defaults
                .build();
    }

    @Bean
    public KStream<String, Account> kStream(StreamsBuilder builder) {

        log.info("Start Kafka Streams application with config: " + properties);
        // configure Serdes
        final Serde<String> stringSerde = Serdes.String();
        final Serde<Account> accountSerde = new SpecificAvroSerde<>();
        accountSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);
        final Serde<Booking> bookingSerde = new SpecificAvroSerde<>();
        bookingSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);
        final Serde<AccountBookings> accountBookingsSerde = new SpecificAvroSerde<>();
        accountBookingsSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);

        // https://kafka.apache.org/30/documentation/streams/developer-guide/dsl-api#ktable-ktable-equi-join
        KStream<String, Account> kStream = builder.stream(properties.getTopicAccounts(), Consumed.with(stringSerde, accountSerde));
        KStream<String, Booking> bookingsKStream = builder.stream(properties.getTopicBookings(), Consumed.with(stringSerde, bookingSerde));

        kStream
                .peek((k, v) -> log.info("Initial record {}", v))
                .leftJoin(countryKTable,
                        (k, v) -> v.getCountry().toString().trim(),
                        (l, r) -> {
                            l.setCountry(r);
                            return l;
                        })

                .peek((k, v) -> log.info("Country added {}", v))
                .selectKey((k, v) -> v.getBookingId().toString(), Named.as("internal-topic")) // rekey
                .peek((k, v) -> log.info("Rekeyed {}", v))
                // .map((k, v) -> KeyValue.pair(v.getBookingId().toString(), v))
                .leftJoin(bookingsKStream,
                        (account, booking) -> map(account, booking),
                        JoinWindows.of(Duration.ofSeconds(60l)),
                        StreamJoined.with(stringSerde, accountSerde, bookingSerde)
                )
                //.leftJoin(bookingsKTable, (account, booking) -> map(account, booking))
                .peek((k, v) -> log.info("Account fully processed {}", v))
                .to(properties.getTopicAccountBookings(), Produced.with(stringSerde, accountBookingsSerde));


        kStream.print(Printed.toSysOut());

        return kStream;

    }
}
