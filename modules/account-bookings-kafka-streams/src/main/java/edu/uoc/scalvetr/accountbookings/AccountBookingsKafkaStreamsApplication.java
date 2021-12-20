package edu.uoc.scalvetr.accountbookings;

import edu.uoc.scalvetr.AccountBookings;
import edu.uoc.scalvetr.Booking;
import edu.uoc.scalvetr.Bookings;
import io.confluent.kafka.serializers.AbstractKafkaSchemaSerDeConfig;
import io.confluent.kafka.streams.serdes.avro.SpecificAvroSerde;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.admin.NewTopic;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.kstream.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Import;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.annotation.EnableKafkaStreams;
import org.springframework.kafka.config.TopicBuilder;

import java.util.ArrayList;
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

    @Autowired
    private KafkaStreamsApplicationProperties properties;

    @Bean
    NewTopic output() {
        return TopicBuilder.name(properties.getOutputTopicName())
                .partitions(properties.getOutputTopicPartitions())
                .replicas(properties.getOutputTopicReplicas())
                .build();
    }

    @Bean
    public KTable<String, event_core_banking_accounts.Value> buildPipeline(StreamsBuilder builder) {

        log.info("Start Kafka Streams application with config: " + properties);
        // configure Serdes
        final Serde<String> stringSerde = Serdes.String();
        final Serde<event_core_banking_accounts.Value> accountSerde = new SpecificAvroSerde<>();
        accountSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);
        final Serde<event_core_banking_bookings.Value> bookingSerde = new SpecificAvroSerde<>();
        bookingSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);
        final Serde<AccountBookings> accountBookingsSerde = new SpecificAvroSerde<>();
        accountBookingsSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);
        final Serde<Bookings> bookingsSerde = new SpecificAvroSerde<>();
        bookingsSerde.configure(Map.of(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG, properties.getSchemaRegistryUrl()), true);

        KTable<String, event_core_banking_accounts.Value> accountKTable = builder.table(properties.getTopicAccounts(), Consumed.with(stringSerde, accountSerde));

        KStream<String, event_core_banking_bookings.Value> bookingKStream = builder.stream(properties.getTopicBookings(), Consumed.with(stringSerde, bookingSerde));

        // aggregate bookings per account id
        KTable<String, Bookings> bookingsKTable = bookingKStream
                .groupByKey().aggregate(Bookings::new, (customerId, booking, bookings) -> {
                            if (bookings.getBookings() == null) {
                                bookings.setBookings(new ArrayList<>());
                            }
                            bookings.getBookings().add(Booking.newBuilder()
                                    .setAccountId(booking.getAccountId())
                                    .setAmount(booking.getAmount())
                                    .setBookingDate(booking.getBookingDate())
                                    .setBookingId(booking.getBookingId())
                                    .setDescription(booking.getDescription())
                                    .setFee(booking.getFee())
                                    .setTaxes(booking.getTaxes())
                                    .setValueDate(booking.getValueDate())
                                    .build());
                            return bookings;
                        },
                        Materialized.as("temp_bookings")
                                .withKeySerde((Serde) stringSerde)
                                .withValueSerde(bookingsSerde));

        //KTable-KTable JOIN to combine account and bookings
        KTable<String, AccountBookings> accountBookingsKTable =
                accountKTable.join(bookingsKTable, (account, bookings) ->
                        AccountBookings.newBuilder()
                                .setAccountId(account.getAccountId())
                                .setBalance(account.getBalance())
                                .setCreationDate(account.getCreationDate())
                                .setIban(account.getIban())
                                .setCancellationDate(account.getCancellationDate())
                                .setCustomerId(account.getCustomerId())
                                .setStatus(account.getStatus())
                                .setBookings(bookings.getBookings()).build());

        accountBookingsKTable.toStream().to(properties.getOutputTopicName(), Produced.with(stringSerde, accountBookingsSerde));

        return accountKTable;
    }

}
