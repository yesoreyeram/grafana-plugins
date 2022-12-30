import { DataSourcePlugin } from '@grafana/data';
import { VercelDS } from './datasource';
import { VercelConfigEditor } from './editors/VercelConfigEditor';
import { VercelQueryEditor } from './editors/VercelQueryEditor';
import { VercelVariablesEditor } from './editors/VercelVariablesEditor';
import type { VercelQuery, VercelConfig, VercelSecureConfig } from './types';

export const plugin = new DataSourcePlugin<VercelDS, VercelQuery, VercelConfig, VercelSecureConfig>(VercelDS)
  .setConfigEditor(VercelConfigEditor)
  .setQueryEditor(VercelQueryEditor)
  .setVariableQueryEditor(VercelVariablesEditor);
