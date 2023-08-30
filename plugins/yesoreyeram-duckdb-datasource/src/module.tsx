import React from 'react';
import { DataSourcePlugin, DataSourceJsonData, DataQuery, DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

type Config = {} & DataSourceJsonData;
type SecureConfig = {};
type Query = {} & DataQuery;
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

const ConfigEditor = () => <>DuckDB Config Editor</>;
const QueryEditor = () => <>DuckDB Query Editor</>;
const VariablesEditor = () => <>DuckDB Variable Editor</>;

export const plugin = new DataSourcePlugin<DataSource, Query, Config, SecureConfig>(DataSource).setConfigEditor(ConfigEditor).setQueryEditor(QueryEditor).setVariableQueryEditor(VariablesEditor);
