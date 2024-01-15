import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface SurrealQuery extends DataQuery {
  queryText: string;
}

export const DEFAULT_QUERY: Partial<SurrealQuery> = {
  queryText: 'SELECT * FROM surreal LIMIT 10',
};

/**
 * These are options configured for each DataSource instance
 */
export interface SurrealDataSourceOptions extends DataSourceJsonData {
  database?: string;
  endpoint?: string;
  namespace?: string;
  scope?: string;
  username?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface SurrealSecureJsonData {
  password?: string;
}
