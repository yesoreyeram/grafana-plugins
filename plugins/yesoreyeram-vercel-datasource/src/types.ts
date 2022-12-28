import { DataSourceJsonData, DataQuery } from '@grafana/data';

export type VercelConfig = {
  apiUrl?: string;
} & DataSourceJsonData;
export type VercelSecureConfig = {
  apiToken: string;
};
export type VercelQuery = {} & DataQuery;
export type VercelVariableQuery = {};
