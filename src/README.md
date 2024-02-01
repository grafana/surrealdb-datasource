<!-- This README file is going to be the one displayed on the Grafana.com website for your plugin -->

# SurrealDB Datasource

The SurrealDB datasource plugin enables you to query and visualize SurrealDB data directly within Grafana, offering seamless integration and exploration of SurrealDB datasets.

## ⚠️ This plugin is currently experimental

This means that while we believe in its potential and are enthusiastic about its development, **we are not yet ready to make a long-term commitment to maintaining it indefinitely**. The plugin is still under active development and may contain bugs. We do not recommend using this plugin in production environments.

## Usage

### Installation

Please refer to our [Data Source Management documentation](https://grafana.com/docs/grafana/latest/administration/data-source-management/) for more information on installing the plugin to an instance of Grafana.

### Configuration

[Add a data source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) by filling in the following fields:

### Basic fields

| Field           | Description                                                                                                       |
| --------------- | ----------------------------------------------------------------------------------------------------------------- |
| Endpoint URL    | The **full** address of the SurrealDB RPC endpoint to connect to, e.g. `ws://localhost:8000/rpc`                  |
| Database name   | The name of the database to connect to.                                                                           |
| Namespace       | The [namespace](https://docs.surrealdb.com/docs/surrealql/statements/define/namespace) to use for the connection. |

### Authentication fields

| Field            | Description                                                                                                     |
| ---------------- | --------------------------------------------------------------------------------------------------------------- |
| Username         | Your SurrealDB username                                                                                         |
| Password         | Your SurrealDB password                                                                                         |
| Scope            | The [scope](https://docs.surrealdb.com/docs/surrealql/statements/define/scope/) to use for the user. (Optional) |

**We strongly recommend that you make your queries with a user account that has read-only access.** This practice not only safeguards your data but also helps maintain system integrity.

### Querying

The query editor allows you to write SurrealQL queries. For more information about writing SurrealQL queries, please refer to [SurrealDB's documentation](https://docs.surrealdb.com/docs/surrealql/overview).

In this version, only a SurrealQL Editor is provided to write queries with. A Query Builder UI is planned for a later version of the plugin.
