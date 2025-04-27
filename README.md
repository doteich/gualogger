# gualogger

**gualogger** is a Go-based OPC UA router designed to efficiently connect to OPC UA servers and forward data changes to various output channels.  It is built with extensibility and cloud-native deployment in mind.

[![Go](https://img.shields.io/badge/Go-1.20-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)  

## Features

* **OPC UA Connectivity:** Connects to OPC UA servers using the `gopcua` library, handling various security policies and modes.
* **Data Routing:** Forwards OPC UA data changes to multiple configurable output channels (exporters).
* **Extensible Architecture:** New output channels can be easily added, allowing for flexible integration with different systems.
* **Cloud Native:** Designed to run in cloud environments, with future Kubernetes integration planned.
* **Configuration via YAML:** Configuration is managed through a YAML file, simplifying setup and modification.
* **Data Transformation:** Basic data type handling and metadata inclusion.
* **Connection Resiliency:** Handles connection retries and keep-alives to ensure robust operation.
* **Certificate Management:** Supports automatic creation or configuration of OPC UA certificates.

## Planned Features

* **Kubernetes Controller:** A Kubernetes controller written in Rust (using `kube-rs`) and Custom Resource Definitions (CRDs) for simplified deployment and management in Kubernetes environments.

## Architecture

The application follows a modular architecture:

* **`main.go`:** The entry point of the application, responsible for loading the configuration and initializing the OPC UA connection and data export manager.
* **`conf.go`:** Handles the loading and parsing of the configuration file (`config.yaml`) using the `viper` library.  Defines the structure of the configuration.
* **`opcua.go`:** Manages the OPC UA client connection, subscription to nodes, and data retrieval.  Includes connection supervision and retry logic.
* **`manager.go`:** The `ExportManager` handles the registration, initialization, and data publishing to the configured output channels (exporters). It also manages metadata associated with Node IDs.
* **`structs.go`:** Defines the data structures used throughout the application, including `Payload` and `Meta`.
* **`handlers/`:** A directory containing the implementations for different output channels (exporters).
    * **`timescale.go`:** Exports data to a TimescaleDB database.
    * **`mqtt.go`:** Exports data to an MQTT broker.
* **`cert.go`:** Handles the creation of self-signed certificates for OPC UA communication.
* **`config_sample.yaml`:** A sample configuration file demonstrating the structure and available options.

## Getting Started

### Prerequisites

* Go 1.20 or later
* An OPC UA server to connect to
* (Optional)  TimescaleDB, MQTT broker, or other services depending on the configured exporters.

### Installation

1.  Clone the repository:

    ```bash
    git clone [<repository_url>](https://github.com/doteich/gualogger.git)
    cd gualogger
    ```

2.  Build the application:

    ```bash
    go build -o gualogger main.go
    ```

### Configuration

1.  Copy `config_sample.yaml` to `config.yaml`:

    ```bash
    cp config_sample.yaml config.yaml
    ```

2.  Edit `config.yaml` to match your environment.  See the "Configuration Details" section below for more information.

#### Sample Config

```yaml
opcua:
  connection:
    endpoint: <OPC UA Server Endpoint IP>
    port: <OPC UA Server Endpoint Port>
    mode: "SignAndEncrypt"   # Possible Entries: 'None', 'Sign', 'SignAndEncrypt'
    policy: 'Basic256Sha256' # Possible Entries: 'None', 'Basic256', 'Basic256Sha256', 'Aes256Sha256RsaPss', 'Aes128Sha256RsaOaep'
    authentication:
      type: 'None'           # Possible Entries: 'None', 'User&Password', 'Certificate'
      credentials:           # Only necessary if type is 'User&Password'
        username: ''
        password: ''
      certificate:           # Only necessary if type is 'Certificate'
        certificate_path: '' # absolute path to certificate file pem encoded
        private_key_path: '' # absolute path to private key file pem encoded
    certificate:             # Only necessary if mode is 'Sign' or 'SignAndEncrypt'
        auto_create: true    # if true, the application will create a self-signed cert on startup, external provided certs are ignored
        certificate_path: '' # absolute path to certificate file used for signing/encryption pem encoded - 
        private_key_path: '' # absolute path to private key file used for signing/encryption pem encoded
    retry_count: 10          # Number of Retries the the connection should retried to the server
  subscription:
    sub_interval: 10         # Subcription Interval in Seconds           
    nodeids:                 # List of Node IDs and associated meta information in key value pairs
      - id: i=2258
        meta:
          - key: mqtt_topic
            value: 'gualogger/topic1'
          - key: name
            value: 'timestamp'
exporters:                   # Map Struct of Exporters - Work in Progress
  timescale-db:
    host: <TimescaleDB Host>
    port: 5432
    username: <TimescaleDB Username>
    password: <TimescaleDB Password>
    database: <TimescaleDB Database>
    table: gualogger
  websocket:
    endpoint: /ws            # Websocket address will be ':{{port}}/{{endpoint}}'
    port: 80
    username: <Websocket Username>
    password: <Websocket Password>
  mqtt:
    protocol: mqtt           # Protocol of the connection string mqtt or mqtts, default is mqtt
    protocol_version: 4      # Protocol version of the connection string, default is 4
    host: <MQTT Broker Host>
    port: 1883
    username: <MQTT Broker Username>
    password: <MQTT Broker Password>
    topic: gualogger
    client_id: gualogger
    qos: 0
    retain: false
    topic_by_nodeid: false   # If true, the topic will be set by the meta information of per nodeid. !Important! the mqtt_topic has to be the first entry per nodeid in the meta object. If false the single topic is set in the mqtt config object
```
### Running the Application

```bash
go run .

