import { DataSourcePlugin } from '@grafana/data';
import { PSConfigEditor, PSQueryEditor, PSVariablesEditor } from './editors';
import { PSDataSource } from './datasource';
import type { PSConfig, PSSecureConfig, PSQuery } from './types';

export const plugin = new DataSourcePlugin<PSDataSource, PSQuery, PSConfig, PSSecureConfig>(PSDataSource)
  .setConfigEditor(PSConfigEditor)
  .setQueryEditor(PSQueryEditor)
  .setVariableQueryEditor(PSVariablesEditor);
