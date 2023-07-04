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
  endpoint?: string;
  username?: string;
  namespace?: string;
  database?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface SurrealSecureJsonData {
  password?: string;
}
