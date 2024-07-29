import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { DataSource } from '../datasource';
import type { QueryEditorProps } from '@grafana/data';
import type { SurrealDataSourceOptions, SurrealQuery } from '../types';

type Props = QueryEditorProps<DataSource, SurrealQuery, SurrealDataSourceOptions>;

export function QueryEditor({ query, onChange }: Props) {
  const onQueryChange = (rawSql: string) => onChange({ ...query, rawSql });

  const { rawSql } = query;

  return (
    <>
      <CodeEditor
        language="sql"
        value={rawSql}
        onBlur={onQueryChange}
        showMiniMap={false}
        showLineNumbers={true}
        height="240px"
      />
    </>
  );
}
