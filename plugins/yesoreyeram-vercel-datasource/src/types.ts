import { DataSourceJsonData, DataQuery } from '@grafana/data';
import type { OpenAPI3Spec } from './types/openapi';

export type VercelConfig = {
  apiUrl?: string;
} & DataSourceJsonData;
export type VercelSecureConfig = {
  apiToken: string;
};
export type KV = { key: string; value: string };
export type VercelQueryType = 'openApi3' | 'raw';
export type VercelQueryBase<T extends VercelQueryType> = { queryType: T } & DataQuery;
export type OpenAPI3Query = {
  type?: 'json' | 'csv' | 'xml' | 'tsv' | 'auto';
  url: string;
  method: 'GET' | 'POST';
  headers?: KV[];
  bodyType?: 'none' | 'raw' | 'form-data' | 'x-www-form-urlencoded';
  body?: string;
  bodyContentType?: string;
  bodyForm?: KV[];
  builder_options: {
    server_url?: string;
    server_variables?: Record<string, string>;
    path?: string;
    path_method?: string;
    path_parameters?: Record<string, string>;
    path_operation_variables?: Record<string, string>;
  };
};
export type VercelQueryOpenApi3 = {
  rootSelector?: string;
} & OpenAPI3Query &
  VercelQueryBase<'openApi3'>;
export type VercelQueryRaw = {} & VercelQueryBase<'raw'>;
export type VercelQuery = VercelQueryOpenApi3 | VercelQueryRaw;
export type VercelVariableQuery = {};

export type GetResourceCallBase<P extends string, Q extends Record<string, any>, R extends unknown> = {
  path: P;
  query?: Q;
  response: R;
};
export type GetResourceCallPing = GetResourceCallBase<'ping', {}, 'pong'>;
export type GetResourceCallOpenAPISpec3 = GetResourceCallBase<'openapi3', {}, OpenAPI3Spec>;
export type GetResourceCall = GetResourceCallPing | GetResourceCallOpenAPISpec3;
