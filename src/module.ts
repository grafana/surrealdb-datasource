import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { SurrealQuery, SurrealDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, SurrealQuery, SurrealDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
