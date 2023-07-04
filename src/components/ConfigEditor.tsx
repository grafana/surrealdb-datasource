import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SurrealDataSourceOptions, SurrealSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<SurrealDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  const onEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      endpoint: event.target.value,
    };

    onOptionsChange({ ...options, jsonData });
  };

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      username: event.target.value,
    };

    onOptionsChange({ ...options, jsonData });
  };

  const onNamespaceChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      namespace: event.target.value,
    };

    onOptionsChange({ ...options, jsonData });
  };

  const onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      database: event.target.value,
    };

    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        password: event.target.value,
      },
    });
  };

  const onResetPassword = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        password: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData: SurrealSecureJsonData = options.secureJsonData ?? {};

  return (
    <div className="gf-form-group">
      <InlineField label="Endpoint" labelWidth={12}>
        <Input
          onChange={onEndpointChange}
          value={jsonData.endpoint || ''}
          placeholder="ws://localhost:8000/rpc"
          width={40}
        />
      </InlineField>
      <InlineField label="Username" labelWidth={12}>
        <Input onChange={onUsernameChange} value={jsonData.username || ''} placeholder="surrealdb" width={40} />
      </InlineField>
      <InlineField label="Password" labelWidth={12}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
          value={secureJsonData.password || ''}
          placeholder="secure json field (backend only)"
          width={40}
          onReset={onResetPassword}
          onChange={onPasswordChange}
        />
      </InlineField>
      <InlineField label="Namespace" labelWidth={12}>
        <Input onChange={onNamespaceChange} value={jsonData.namespace || ''} placeholder="my-department" width={40} />
      </InlineField>
      <InlineField label="Database" labelWidth={12}>
        <Input onChange={onDatabaseChange} value={jsonData.database || ''} placeholder="grafana-data" width={40} />
      </InlineField>
    </div>
  );
}
