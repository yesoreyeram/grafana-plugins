import React, { useState } from 'react';
import { InlineLabel, Button, Input, RadioButtonGroup, useTheme2 } from '@grafana/ui';
import { DataSourcePlugin, DataSourcePluginOptionsEditorProps, QueryEditorProps, DataQuery, DataSourceJsonData, DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

type Config = {} & DataSourceJsonData;
type SecureConfig = { api_token: string };

type QueryType = 'monitors' | 'maintenance-windows';
type BaseQuery<T extends QueryType> = { queryType: T } & DataQuery;
type MonitorsQuery = {} & BaseQuery<'monitors'>;
type MaintenanceWindowsQuery = {} & BaseQuery<'maintenance-windows'>;
type Query = MonitorsQuery | MaintenanceWindowsQuery;

type VariableQuery = {};

class DataSource extends DataSourceWithBackend<Query, Config> {
  constructor(instanceSettings: DataSourceInstanceSettings<Config>) {
    super(instanceSettings);
    this.annotations = {};
  }
  filterQuery(query: Query): boolean {
    return !query.hide;
  }
  metricFindQuery(query: VariableQuery, options: unknown): Promise<MetricFindValue[]> {
    return new Promise((resolve) => resolve([]));
  }
}

const ConfigEditor = (props: DataSourcePluginOptionsEditorProps<Config, SecureConfig>) => {
  const { options, onOptionsChange } = props;
  const { secureJsonFields } = options;
  const { secureJsonData = { api_token: '' } } = options;
  const [api_token, set_api_token] = useState('');
  const theme = useTheme2();
  const onSecureOptionChange = <Key extends keyof SecureConfig, Value extends SecureConfig[Key]>(option: Key, value: Value, set: boolean) => {
    onOptionsChange({
      ...options,
      secureJsonData: { ...secureJsonData, [option]: value },
      secureJsonFields: { ...secureJsonFields, [option]: set },
    });
  };
  return (
    <>
      <div className="gf-form">
        <InlineLabel tooltip={'hyperping API Token'} width={20}>
          API Token
        </InlineLabel>
        {secureJsonFields?.api_token ? (
          <>
            <Input type="text" value="Configured" disabled={true} width={40}></Input>
            <Button
              style={{ marginInlineStart: theme.spacing(0.5) }}
              variant="secondary"
              className="reset-button"
              onClick={() => {
                set_api_token('');
                onSecureOptionChange('api_token', api_token, false);
              }}
            >
              Reset
            </Button>
          </>
        ) : (
          <Input
            width={40}
            type="password"
            placeholder="hyperping API Token"
            value={api_token || ''}
            onChange={(e) => set_api_token(e.currentTarget.value)}
            onBlur={() => onSecureOptionChange('api_token', api_token, true)}
          />
        )}
      </div>
    </>
  );
};

const QueryEditor = (props: QueryEditorProps<DataSource, Query, Config>) => {
  const { query, onChange, onRunQuery } = props;
  return (
    <>
      <div className="gf-form">
        <InlineLabel width={20}>Query Type</InlineLabel>
        <RadioButtonGroup<QueryType>
          options={[
            { value: 'monitors', label: 'Monitors' },
            { value: 'maintenance-windows', label: 'Maintenance windows' },
          ]}
          value={query.queryType || 'monitors'}
          onChange={(e) => {
            onChange({ ...query, queryType: e || 'monitors' });
            onRunQuery();
          }}
        ></RadioButtonGroup>
      </div>
    </>
  );
};

const QueryEditorHelp = () => {
  return <></>;
};

const VariablesEditor = () => {
  return <>HyperPing Variable Editor</>;
};

export const plugin = new DataSourcePlugin<DataSource, Query, Config, SecureConfig>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setQueryEditorHelp(QueryEditorHelp)
  .setVariableQueryEditor(VariablesEditor);
