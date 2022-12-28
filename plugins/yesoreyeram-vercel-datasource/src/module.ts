import { DataSourcePlugin } from '@grafana/data';
import { VercelDS } from './datasource';
import { VercelConfigEditor, VercelQueryEditor, VercelVariablesEditor } from './editors';
import type { VercelQuery, VercelConfig, VercelSecureConfig } from './types';

export const plugin = new DataSourcePlugin<VercelDS, VercelQuery, VercelConfig, VercelSecureConfig>(VercelDS)
  .setConfigEditor(VercelConfigEditor)
  .setQueryEditor(VercelQueryEditor)
  .setVariableQueryEditor(VercelVariablesEditor);
