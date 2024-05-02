# Changelog

## 0.2.4 (2024-05-02)

- ko.yaml: import embedded timezone database into application
## 0.2.3 (2024-05-02)

- (sinks/gotify) hotfix for title

## 0.2.2 (2024-05-02)

- (sinks/gotify) add template support for gotify title, example:
    add in config:
    ```yml
    ...
    sinks:
      gotify:
        params:
          title: '{{ now "15:04" }}: report from machine'
    ...
    ```
    This will add hour and minute in title of gotify message.

## 0.2.1 (2024-01-04)

- (sinks/windy) fix sink when station id is 0

## 0.2.0 (2023-11-25)

- (dispatcher/event) fix double event trigger when multiple sinks are assigned
- (collector/mqtt) subscribe after reconnect

## 0.1.3 (2023-11-24)

- (sinks/gotify) fixed priority was not being set in message from params
- (README) update info about environment variables and docker deployment

## 0.1.2 (2023-11-23)

- added 32 bit arm containers

## 0.1.1 (2023-11-23)

- added multi platform container images builds to CI config

## 0.1.0 (2023-11-23)

- first versioned release
- all basic functionality is working
- not tested very well
