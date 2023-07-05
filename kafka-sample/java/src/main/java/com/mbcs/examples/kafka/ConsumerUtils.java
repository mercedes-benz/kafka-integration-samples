package com.mbcs.examples.kafka;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.FileReader;
import java.io.IOException;
import java.time.Duration;
import java.util.List;
import java.util.Properties;

public class ConsumerUtils {
    private static final Logger log = LoggerFactory.getLogger(ConsumerUtils.class);

    private static Properties loadConfig(String filePath) throws IOException {
        final Properties cfg = new Properties();
        try (FileReader configFile = new FileReader(filePath)) {
            cfg.load(configFile);
            return cfg;
        }
    }

    /**
     * Convenience method to open a kafka consumer which consumes from the specified topic.
     *
     * @param topicName  topic name, from which the consumer should read from.
     * @param configPath path to config file for the consumer
     * @throws IOException when config file can't be opened
     */
    public static void consumeFromTopic(String topicName, String configPath) throws IOException {
        Properties props = loadConfig(configPath);
        try (final KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props)) {
            consumer.subscribe(List.of(topicName));
            while (true) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    String key = record.key();
                    String value = record.value();
                    log.debug("received event from topic {} with key {}", topicName, key);
                    System.out.println(value);
                }
            }
        }
    }
}