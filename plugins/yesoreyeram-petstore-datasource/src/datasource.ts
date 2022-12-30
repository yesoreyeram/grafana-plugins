import { DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import type { PSConfig, PSQuery, PSVariableQuery, GetResourceCall, GetResourceCallOpenAPISpec3, OpenAPI3Spec } from './types';

export class PSDataSource extends DataSourceWithBackend<PSQuery, PSConfig> {
  constructor(instanceSettings: DataSourceInstanceSettings<PSConfig>) {
    super(instanceSettings);
    this.annotations = {};
  }
  filterQuery(query: PSQuery): boolean {
    return !query.hide;
  }
  metricFindQuery(query: PSVariableQuery, options: unknown): Promise<MetricFindValue[]> {
    return new Promise((resolve) => resolve([]));
  }
  getResource<O extends GetResourceCall>(path: O['path'], params?: O['query']): Promise<O['response']> {
    return super.getResource(path, params);
  }
  getOpenAPI3Spec = (): Promise<OpenAPI3Spec> => {
    return this.getResource<GetResourceCallOpenAPISpec3>('openapi3');
  };
}
