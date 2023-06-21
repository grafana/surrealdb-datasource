import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SurrealQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const DEFAULT_QUERY: Partial<SurrealQuery> = {
  constant: 6.5,
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
