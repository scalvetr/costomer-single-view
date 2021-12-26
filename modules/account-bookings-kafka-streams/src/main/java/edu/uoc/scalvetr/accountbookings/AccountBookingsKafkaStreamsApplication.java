package edu.uoc.scalvetr.accountbookings;

import edu.uoc.scalvetr.AccountBookings;
import edu.uoc.scalvetr.Booking;
import edu.uoc.scalvetr.Bookings;
import io.confluent.kafka.schemaregistry.client.rest.entities.SchemaString;
import io.confluent.kafka.serializers.AbstractKafkaSchemaSerDeConfig;
import io.confluent.kafka.serializers.KafkaAvroSerializer;
import io.confluent.kafka.serializers.context.NullContextNameStrategy;
import io.confluent.kafka.serializers.subject.TopicNameStrategy;
import io.confluent.kafka.streams.serdes.avro.SpecificAvroSerde;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.admin.NewTopic;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.Topology;
import org.apache.kafka.streams.errors.DefaultProductionExceptionHandler;
import org.apache.kafka.streams.errors.LogAndFailExceptionHandler;
import org.apache.kafka.streams.kstream.*;
import org.apache.kafka.streams.processor.FailOnInvalidTimestamp;
import org.apache.kafka.streams.processor.internals.StreamsPartitionAssignor;
import org.apache.kafka.streams.processor.internals.assignment.FallbackPriorTaskAssignor;
import org.apache.kafka.streams.processor.internals.assignment.HighAvailabilityTaskAssignor;
import org.apache.kafka.streams.processor.internals.assignment.StickyTaskAssignor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Import;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.annotation.EnableKafkaStreams;
import org.springframework.kafka.config.TopicBuilder;
import org.springframework.messaging.support.ErrorMessage;
import org.springframework.nativex.hint.NativeHint;
import org.springframework.nativex.hint.TypeHint;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@SpringBootApplication
@EnableConfigurationProperties
@EnableKafka
@EnableKafkaStreams
@Import(KafkaStreamsApplicationProperties.class)
@RequiredArgsConstructor
@Slf4j
// hints for spring native
@NativeHint(
        types = @TypeHint(types = {
                StreamsPartitionAssignor.class,
                DefaultProductionExceptionHandler.class,
                FailOnInvalidTimestamp.class,
                HighAvailabilityTaskAssignor.class,
                StickyTaskAssignor.class,
                FallbackPriorTaskAssignor.class,
                TopicNameStrategy.class,
                LogAndFailExceptionHandler.class,
                SpecificAvroSerde.class,
                KafkaAvroSerializer.class,
                ErrorMessage.class,
                SchemaString.class,
                NullContextNameStrategy.class
        })
)
@TypeHint(types = Serdes.class,
        typeNames = {
                "org.apache.kafka.common.serialization.Serdes$StringSerde",
                "org.apache.kafka.common.serialization.Serdes$ByteArraySerde"
        })
public class AccountBookingsKafkaStreamsApplication {

    public static void main(String[] args) {
        SpringApplication.run(AccountBookingsKafkaStreamsApplication.class, args);
    }

    @Autowired
    private KafkaStreamsApplicationProperties properties;

    @Bean
    NewTopic output() {
        // Controls the partitions and replicas of the output topic
        return TopicBuilder.name(properties.getOutputTopicName())
                .partitions(properties.getOutputTopicPartitions())
                .replicas(properties.getOutputTopicReplicas())
                .build();
    }

    @Bean
    public Topology buildPipeline(StreamsBuilder builder) {
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

        // aggregate bookings in the bookings stream by account id
        KStream<String, event_core_banking_bookings.Value> bookingKStream = builder.stream(properties.getTopicBookings(), Consumed.with(stringSerde, bookingSerde));
        KTable<String, Bookings> bookingsKTable = bookingKStream.groupByKey() // booking's KStream key is account_id
                .aggregate(Bookings::new, (customerId, booking, bookings) -> { // aggregate bookings in a new structure called Bookings
                            if (bookings.getBookings() == null) {
                                bookings.setBookings(new ArrayList<>());
                            }
                            bookings.getBookings().add(buildBooking(booking));
                            return bookings;
                        },
                        Materialized.as("temp_bookings") // materialize for fault tolerance
                                .withKeySerde((Serde) stringSerde)
                                .withValueSerde(bookingsSerde));

        //KTable-KTable JOIN to combine account and bookings
        KTable<String, event_core_banking_accounts.Value> accountKTable = builder.table(properties.getTopicAccounts(), Consumed.with(stringSerde, accountSerde));
        KTable<String, AccountBookings> accountBookingsKTable =
                accountKTable.join(bookingsKTable, (account, bookings) -> // Join accounts with the booking aggregation previously produced
                        buildAccountBooking(account, bookings.getBookings())); // build the AccountBookings

        // convert to stream and produce to the output topic
        accountBookingsKTable.toStream().to(properties.getOutputTopicName(), Produced.with(stringSerde, accountBookingsSerde));
        return builder.build();
    }

    private AccountBookings buildAccountBooking(event_core_banking_accounts.Value account, List<Booking> bookings) {
        return AccountBookings.newBuilder()
                .setAccountId(account.getAccountId())
                .setBalance(account.getBalance())
                .setCreationDate(account.getCreationDate())
                .setIban(account.getIban())
                .setCancellationDate(account.getCancellationDate())
                .setCustomerId(account.getCustomerId())
                .setStatus(account.getStatus())
                .setBookings(bookings).build();
    }

    private Booking buildBooking(event_core_banking_bookings.Value booking) {
        return Booking.newBuilder()
                .setAccountId(booking.getAccountId())
                .setAmount(booking.getAmount())
                .setBookingDate(booking.getBookingDate())
                .setBookingId(booking.getBookingId())
                .setDescription(booking.getDescription())
                .setFee(booking.getFee())
                .setTaxes(booking.getTaxes())
                .setValueDate(booking.getValueDate())
                .build();
    }

}
