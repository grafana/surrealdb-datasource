import React, { ChangeEvent } from 'react';
import { Alert, Divider, Field, Input, SecretInput, Stack, TextLink } from '@grafana/ui';
import { DataSourceDescription, ConfigSection } from '@grafana/plugin-ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import type { SurrealDataSourceOptions, SurrealSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<SurrealDataSourceOptions> {}

export function ConfigEditor({ onOptionsChange, options }: Props) {
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
      <Alert title="SurrealDB v2.0 compatibility" severity="warning">
        <Stack direction="column">
          <div>
            The Grafana SurrealDB datasource currently does not support SurrealDB v2.0. Please ensure you are using a
            compatible version of SurrealDB (v1.x) for full functionality. Follow the GitHub issue{' '}
            <TextLink href="https://github.com/grafana/surrealdb-datasource/issues/441" external inline>
              here
            </TextLink>{' '}
            for updates on compatibility.
          </div>
        </Stack>
      </Alert>

      <Alert title="This datasource is currently experimental" severity="warning">
        <Stack direction="column">
          <div>
            This means that you might encounter unexpected behavior, bugs, or limitations while using this datasource.
            We strongly advise exercising caution and understanding the potential risks associated with using
            experimental software.
          </div>
          <div>
            Found a bug? Have a suggestion? Open an issue on{' '}
            <TextLink href="https://github.com/grafana/surrealdb-datasource/issues" external inline>
              Github
            </TextLink>
            !
          </div>
        </Stack>
      </Alert>

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
