import { DataSourceInstanceSettings, CoreApp } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { SurrealQuery, SurrealDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<SurrealQuery, SurrealDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SurrealDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<SurrealQuery> {
    return DEFAULT_QUERY
  }
}
