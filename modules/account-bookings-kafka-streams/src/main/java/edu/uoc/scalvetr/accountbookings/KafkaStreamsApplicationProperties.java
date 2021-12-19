package edu.uoc.scalvetr.accountbookings;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties("config")
@Data
public class KafkaStreamsApplicationProperties {

    private String kafkaBrokers;
    private String kafkaUsername;
    private String kafkaPassword;
    private String applicationId;
    private String schemaRegistryUrl;
    private String topicAccounts;
    private String topicBookings;
    private String topicAccountBookings;

}
