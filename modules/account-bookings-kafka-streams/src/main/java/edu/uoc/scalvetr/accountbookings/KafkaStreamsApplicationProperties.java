package edu.uoc.scalvetr.accountbookings;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties("config")
@Data
public class KafkaStreamsApplicationProperties {

    private String kafkaBrokers;
    private String applicationId;
    private String schemaRegistryUrl;
    private String topicAccounts;
    private String topicBookings;
    private String outputTopicName;
    private int outputTopicPartitions = 1;
    private int outputTopicReplicas = 1;

}
