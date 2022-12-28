import { DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import type { VercelQuery, VercelConfig, VercelVariableQuery } from './types';

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
}
