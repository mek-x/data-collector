# data-collector

[![License
MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)
[![Go Report
Card](https://goreportcard.com/badge/gitlab.com/mek_x/data-collector)](https://goreportcard.com/report/gitlab.com/mek_x/data-collector)

This is a program used to gather data from various sources (like MQTT topics) and push the data to different sinks (e.g.
to file, Gotify notification service, IOTPlotter, database). Dynamically configurable and intended to be deployed as a
containerized service.

## Overview

Overall architecture:

```mermaid
flowchart LR
    A1[Collector A] --> B1[Parser]
    A2[Collector B] --> B2[Parser]
    B1 & B2 --> C(DataStore)
    C --> D1[Dispatcher A]
    C --> D2[Dispatcher B]
    D1 --> S1[Sink A]
    D1 --> S2[Sink B]
    D2 --> S3[Sink C]
    D2 --> S4[Sink D]
```

Legend:

1. Collectors - collects data from configured sources
    - MQTT - subscribes to configured topics and gathers the data
    - File - gathers the data from files
    - Shell - gets data from shell
2. Parsers - convert data from source format into the intermediate value to be used by dispatcher/sink
    - jsonpath - uses jsonpath to parse data
    - regex - uses regular expression to parse data
3. DataStore - main data structure, where data is held between publishing by collector and sending by dispatcher/sink
4. Dispatchers - schedules which chunk of data will be dispatched
    - cron: uses cron-like interface to schedule when which data should be dispatched
    - event: event triggered (e.g. on value changed, on threshold)
5. Sinks - dispatches the scheduled data to the sink
    - REST-based sinks - for sinks using REST interface (e.g. gotify)
    - file sink - for sinks using file interface (e.g. csv file writer)

## Example configuration

Configuration is in YAML file. Example is [here](./configs/example.yaml)

### Support for evaluating expressions

Used in:
- `data.<name>.conv` - to convert variable, e.g. multiply the value by 10 (`val * 10`)
- `dispatchers.sink[].spec` - when `type` is `expr` it can be used for templating, e.g. to generate json (`toJSON(val)`)

Syntax for the evaluation is documented here: [Expr Language Definition](https://expr.medv.io/docs/Language-Definition)

## Deployment

This project uses `ko` as a build system.

## Development

### Adding new modules

Adding additional collectors, parsers, dispatchers and sinks is pretty much self-contained. Only necessary things is to
add the module code itself and import it in `internal/modules.modules.go`. Take a look at the other modules as an
example, like [file collector](pkg/collector/file/file.go).

## Todo list

1. Collectors:
    - [x] mqtt
    - [x] file
    - [ ] shell
2. Parsers:
    - [x] add expression eval support
    - [x] jsonpath
    - [ ] regex
3. Dispatchers:
    - [x] cron
    - [ ] event
4. Sinks:
    - [x] support for templating output
    - [ ] generic REST (template/json)
    - [ ] gotify (template)
    - [ ] iotplotter (json)
    - [ ] file (csv)
    - [x] stdout (for debugging or template)
5. Configuration:
    - [x] yaml support
    - [x] replacing env variables
        - [ ] testing needed
6. Others:
    - [ ] CI

## Dependencies

- [jsonpath](https://github.com/PaesslerAG/jsonpath)
- [Paho MQTT](https://github.com/eclipse/paho.mqtt.golang)
- [gocron](https://github.com/go-co-op/gocron)
- [go-yaml](https://github.com/goccy/go-yaml)
- [expr](https://github.com/antonmedv/expr)
- [mapstructure](https://github.com/mitchellh/mapstructure)
