//#region Imports
import React, { useState } from 'react';
import { DataSourcePlugin, DataSourceJsonData, DataQuery, DataSourceInstanceSettings, DataSourcePluginOptionsEditorProps, QueryEditorProps, MetricFindValue, ScopedVars } from '@grafana/data';
import { InlineFormLabel, Input, Select, Button } from '@grafana/ui';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { SecretMessage } from '@yesoreyeram/grafana-plugins-ui';
import { echo } from '@yesoreyeram/grafana-plugins-utils';
//#endregion
//#region Types - Configuration
type NameTransformMode = 'none' | 'upper_case' | 'lower_case';
type Config = { nameTransformMode: NameTransformMode } & DataSourceJsonData;
type SecureConfig = { apiToken: string };
//#endregion
//#region Types - Query
type QueryType = 'greet';
type QueryBase<QT extends QueryType> = { queryType: QT } & DataQuery;
type GreetQuery = { greeting?: string; username: string } & QueryBase<'greet'>;
export type Query = GreetQuery;
//#endregion
//#region Types - Variable Query
type VariableQueryType = 'greeting-list';
type VariableQueryBase<QT extends VariableQueryType> = { queryType: QT };
type GreetingsListVariableQuery = {} & VariableQueryBase<'greeting-list'>;
type VariableQuery = GreetingsListVariableQuery;
//#endregion
//#region Types - Resource calls
type GetResourceCallPath = 'greeting-list';
type GetResourceCallBase<P extends GetResourceCallPath, Q extends Record<string, any>, R extends unknown> = {
  path: P;
  query?: Q;
  response: R;
};
type GetResourceCallGreetingsList = GetResourceCallBase<'greeting-list', {}, string[]>;
type GetResourceCall = GetResourceCallGreetingsList;
//#endregion
//#region Migration and defaults
const applyDefaultsToQuery = (source_query: Partial<Query> = {}, instanceSettings: DataSourceInstanceSettings<Config>): Query => {
  const query: Query = { ...source_query } as Query;
  if (!query.queryType) {
    query.queryType = 'greet';
  }
  return query;
};
//#endregion
//#region Interpolation
const applyTemplateVariables = (query: Query, scopedVars: ScopedVars): Query => {
  switch (query.queryType) {
    case 'greet':
      query = {
        ...query,
        greeting: getTemplateSrv().replace(query.greeting || '', scopedVars),
        username: getTemplateSrv().replace(query.username || '', scopedVars),
      };
  }
  return query;
};
//#endregion
//#region DataSource
export class DataSource extends DataSourceWithBackend<Query, Config> {
  constructor(private instanceSettings: DataSourceInstanceSettings<Config>) {
    super(instanceSettings);
    this.annotations = {};
  }
  filterQuery(query: Query): boolean {
    return !query.hide;
  }
  metricFindQuery(query: VariableQuery, options: unknown): Promise<MetricFindValue[]> {
    return new Promise((resolve, reject) => {
      switch (query.queryType) {
        case 'greeting-list':
        default:
          this.getGreetingList()
            .then((result) => {
              resolve(result.map((r) => ({ text: r, value: r })));
            })
            .catch(reject);
      }
    });
  }
  interpolateVariablesInQueries(queries: Query[], scopedVars: ScopedVars = {}): Query[] {
    return queries.map((query) => {
      return this.applyTemplateVariables(query, scopedVars);
    });
  }
  applyTemplateVariables(query: Query, scopedVars: ScopedVars = {}): Query {
    let newQuery = applyDefaultsToQuery(query || {}, this.instanceSettings);
    return applyTemplateVariables(newQuery, scopedVars);
  }
  getResource<O extends GetResourceCall>(path: O['path'], params?: O['query']): Promise<O['response']> {
    return super.getResource(path, params);
  }
  private getGreetingList(): Promise<string[]> {
    return this.getResource<GetResourceCallGreetingsList>('greeting-list');
  }
}
//#endregion
//#region Editors
const ConfigEditor = (props: DataSourcePluginOptionsEditorProps<Config, SecureConfig>) => {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonFields } = options;
  const { secureJsonData = { apiToken: '' } } = options;
  const [apiToken, setApiToken] = useState('');
  const onOptionChange = <Key extends keyof Config, Value extends Config[Key]>(option: Key, value: Value) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, [option]: value },
    });
  };
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
        <InlineFormLabel width={12}>Name Transform Mode</InlineFormLabel>
        <Select<NameTransformMode>
          value={jsonData.nameTransformMode || 'none'}
          options={[
            { value: 'none', label: 'None' },
            { value: 'lower_case', label: 'Lower Case' },
            { value: 'upper_case', label: 'Upper Case' },
          ]}
          onChange={(e) => onOptionChange('nameTransformMode', e.value!)}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel tooltip={'API Token. Use hello as token'} width={12}>
          API Token
        </InlineFormLabel>
        {secureJsonFields?.apiToken ? (
          <>
            <Input type="text" value="Configured" disabled={true}></Input>
            <Button
              style={{ marginInlineStart: '10px' }}
              variant="secondary"
              className="reset-button"
              onClick={() => {
                setApiToken('');
                onSecureOptionChange('apiToken', apiToken, false);
              }}
            >
              Reset
            </Button>
          </>
        ) : (
          <Input
            type="password"
            placeholder="API Token. Use hello as api token"
            value={apiToken || ''}
            onChange={(e) => setApiToken(e.currentTarget.value)}
            onBlur={() => onSecureOptionChange('apiToken', apiToken, true)}
          />
        )}
      </div>
      <SecretMessage message={echo('Hello datasource plugin')} />
    </>
  );
};
export const QueryEditor = (props: QueryEditorProps<DataSource, Query, Config>) => {
  const { query, onChange, onRunQuery } = props;
  const [greeting, setGreeting] = useState(query.greeting || 'Hello');
  const [username, setUsername] = useState(query.username || 'Grafana User');
  return (
    <>
      <div className="gf-form">
        <InlineFormLabel>Query Type</InlineFormLabel>
        <Select<QueryType> options={[{ value: 'greet', label: 'Greeting' }]} value={query.queryType || 'greet'} onChange={(e) => onChange({ ...query, queryType: e.value! })} />
      </div>
      <div className="gf-form">
        <InlineFormLabel tooltip={'Enter how the user want to be greeted'}>Greeting</InlineFormLabel>
        <Input
          value={greeting}
          onChange={(e) => setGreeting(e.currentTarget.value)}
          onBlur={() => {
            onChange({ ...query, greeting });
            onRunQuery();
          }}
        />
      </div>{' '}
      <div className="gf-form">
        <InlineFormLabel tooltip={'Enter your full name here'}>User Name</InlineFormLabel>
        <Input
          value={username}
          onChange={(e) => setUsername(e.currentTarget.value)}
          onBlur={() => {
            onChange({ ...query, username });
            onRunQuery();
          }}
        />
      </div>
    </>
  );
};
const VariablesEditor = (props: { query: VariableQuery; onChange: (query: VariableQuery, definition: string) => void; datasource: DataSource }) => {
  const { query, onChange } = props;
  return (
    <>
      <div className="gf-form">
        <InlineFormLabel>Query Type</InlineFormLabel>
        <Select<typeof query.queryType>
          options={[{ value: 'greeting-list', label: 'Greeting List' }]}
          value={query.queryType}
          onChange={(e) => onChange({ ...query, queryType: e.value! }, query.queryType)}
        />
      </div>
    </>
  );
};
//#endregion
//#region Module
export const plugin = new DataSourcePlugin<DataSource, Query, Config, SecureConfig>(DataSource).setConfigEditor(ConfigEditor).setQueryEditor(QueryEditor).setVariableQueryEditor(VariablesEditor);
//#endregion
