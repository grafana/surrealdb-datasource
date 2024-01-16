import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { DataSource } from '../datasource';
import type { QueryEditorProps } from '@grafana/data';
import type { SurrealDataSourceOptions, SurrealQuery } from '../types';

type Props = QueryEditorProps<DataSource, SurrealQuery, SurrealDataSourceOptions>;

export function QueryEditor({ query, onChange }: Props) {
  const onQueryTextChange = (queryText: string) => onChange({ ...query, queryText });

  const { queryText } = query;

  return (
    <>
      <CodeEditor
        language="sql"
        value={queryText}
        onBlur={onQueryTextChange}
        showMiniMap={false}
        showLineNumbers={true}
        height="240px"
      />
    </>
  );
}
