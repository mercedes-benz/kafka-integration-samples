package com.mbcs.examples.kafka;

import org.apache.commons.cli.*;

public class KafkaSampleCliParser {
    private final CommandLineParser parser = new DefaultParser();
    private final Options options = new Options();

    public CommandLine parse(String[] args) throws ParseException {

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
        cmd = parser.parse(options, args);
        return cmd;
    }

    public void printHelp() {
        HelpFormatter hf = new HelpFormatter();
        hf.printHelp("java -jar <executable .jar path>", options, true);
    }
}
