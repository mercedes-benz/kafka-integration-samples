package com.mbcs.examples.kafka;

import org.apache.commons.cli.*;

import java.io.IOException;

public class Main {
    public static void main(String[] args) {
        CommandLineParser parser = new DefaultParser();

        Options options = new Options();
        options.addOption(Option.builder("t")
                .longOpt("topic")
                .desc("topic to subscribe to")
                .required()
                .hasArg()
                .argName("topic name")
                .type(String.class)
                .build());
        options.addOption(Option.builder("c")
                .longOpt("config")
                .desc("path to kafka consumer config file")
                .hasArg()
                .argName("config file path")
                .type(String.class)
                .build());

        CommandLine cmd;
        try {
            cmd = parser.parse(options, args);
        } catch (ParseException e) {
            printHelp(options);
            return;
        }

        String topic = cmd.getOptionValue("t");
        String consumerConfig = cmd.getOptionValue("c", "consumer.properties");
        try {
            ConsumerUtils.consumeFromTopic(topic, consumerConfig);
        } catch (IOException e) {
            System.err.println("Consumer config file '" + consumerConfig + "' could not be opened.");
            printHelp(options);
        }
    }

    private static void printHelp(Options options) {
        HelpFormatter hf = new HelpFormatter();
        hf.printHelp("java -jar <executable .jar path>", options, true);
    }
}