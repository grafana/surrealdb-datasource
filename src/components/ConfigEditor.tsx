import React, { ChangeEvent } from 'react';
import { Divider, Field, Input, SecretInput } from '@grafana/ui';
import { DataSourceDescription, ConfigSection } from '@grafana/experimental';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import type { SurrealDataSourceOptions, SurrealSecureJsonData } from '../types';

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

  const onScopeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      scope: event.target.value,
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
    <>
      <DataSourceDescription
        dataSourceName="SurrealDB"
        docsLink="https://grafana.com/grafana/plugins/surrealdb-datasource/"
        hasRequiredFields
      />
      <Divider />
      <ConfigSection title="Server">
        <Field
          required
          label={'Endpoint URL'}
          description={'The address of the SurrealDB server to connect to.'}
          invalid={!jsonData.endpoint}
          error={'Endpoint URL is required'}
        >
          <Input
            name="endpoint"
            width={40}
            value={jsonData.endpoint || ''}
            onChange={onEndpointChange}
            label={'Endpoint URL'}
            aria-label={'Endpoint URL'}
            placeholder={'ws://localhost:8000/rpc'}
          />
        </Field>
        <Field
          required
          label={'Database name'}
          description={'The name of the database to connect to.'}
          invalid={!jsonData.database}
          error={'Database name is required'}
        >
          <Input
            name="port"
            width={40}
            value={jsonData.database || ''}
            onChange={onDatabaseChange}
            label={'Database name'}
            aria-label={'Database name'}
            placeholder={'Database name'}
          />
        </Field>
        <Field
          required
          label={'Namespace'}
          description={'The namespace to use for the connection.'}
          invalid={!jsonData.namespace}
          error={'Namespace is required'}
        >
          <Input
            name="namespace"
            width={40}
            value={jsonData.namespace || ''}
            onChange={onNamespaceChange}
            label={'Namespace'}
            aria-label={'Namespace'}
            placeholder={'Namespace'}
          />
        </Field>
      </ConfigSection>
      <Divider />
      <ConfigSection title="Authentication">
        <Field
          required
          label={'Username'}
          description={'The username to use for the connection.'}
          invalid={!jsonData.username}
          error={'Username is required'}
        >
          <Input
            name="username"
            width={40}
            value={jsonData.username || ''}
            onChange={onUsernameChange}
            label={'Username'}
            aria-label={'Username'}
            placeholder={'Username'}
          />
        </Field>
        <Field required label={'Password'} description={'The password to use for the connection.'}>
          <SecretInput
            name="pwd"
            width={40}
            label={'Password'}
            aria-label={'Password'}
            placeholder={'Password'}
            value={secureJsonData.password || ''}
            isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
            onReset={onResetPassword}
            onChange={onPasswordChange}
          />
        </Field>
        <Field label={'Scope'} description={'The scope to use for the connection.'}>
          <Input
            name="scope"
            width={40}
            value={jsonData.scope || ''}
            onChange={onScopeChange}
            label={'Scope'}
            aria-label={'Scope'}
            placeholder={'Scope'}
          />
        </Field>
      </ConfigSection>
      <Divider />
      <ConfigSection
        title="Additional settings"
        description="Additional settings are optional settings that can be configured for more control over your data source."
        isCollapsible
        isInitiallyOpen={true}
      >
        Not available in this version.
      </ConfigSection>
    </>
  );
}
