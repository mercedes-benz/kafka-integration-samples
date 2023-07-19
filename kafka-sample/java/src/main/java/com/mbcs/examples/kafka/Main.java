package com.mbcs.examples.kafka;


import org.apache.commons.cli.CommandLine;

import java.io.IOException;

public class Main {

    public static void main(String[] args) {
        ConsumerUtils consumer = new ConsumerUtils();
        KafkaSampleCliParser parser = new KafkaSampleCliParser();

        String consumerConfig = null;
        try {
            CommandLine cmd = parser.parse(args);
            String topic = cmd.getOptionValue("t");
            consumerConfig = cmd.getOptionValue("c", "consumer.properties");

            consumer.consumeFromTopic(topic, consumerConfig);

        } catch (IOException e) {
            System.err.println("Consumer config file '" + consumerConfig + "' could not be opened.");
            parser.printHelp();
        } catch (Exception e) {
            parser.printHelp();
        }
    }

}