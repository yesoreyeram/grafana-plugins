import { DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import type { OpenAPI3Spec } from './types/openapi';
import type { VercelQuery, VercelConfig, VercelVariableQuery, GetResourceCall, GetResourceCallPing, GetResourceCallOpenAPISpec3 } from './types';

export class VercelDS extends DataSourceWithBackend<VercelQuery, VercelConfig> {
  constructor(instanceSettings: DataSourceInstanceSettings<VercelConfig>) {
    super(instanceSettings);
    this.annotations = {};
  }
  filterQuery(query: VercelQuery): boolean {
    return !query.hide;
  }
  metricFindQuery(query: VercelVariableQuery, options: unknown): Promise<MetricFindValue[]> {
    return new Promise((resolve) => resolve([]));
  }
  getResource<O extends GetResourceCall>(path: O['path'], params?: O['query']): Promise<O['response']> {
    return super.getResource(path, params);
  }
  getResourcePing = (): Promise<'pong'> => {
    return this.getResource<GetResourceCallPing>('ping');
  };
  getOpenAPI3Spec = (): Promise<OpenAPI3Spec> => {
    return this.getResource<GetResourceCallOpenAPISpec3>('openapi3');
  };
}
