import React, { useState, useEffect } from 'react';
import { DataSourcePluginOptionsEditorProps, QueryEditorProps } from '@grafana/data';
import { OpenAPI3Editor } from './editors/OpenApi3Editor';
import { PSDataSource } from './datasource';
import type { PSConfig, PSSecureConfig, PSQuery, OpenAPI3Query, PSVariableQuery, OpenAPI3Spec } from './types';

type PSQueryEditorProps = QueryEditorProps<PSDataSource, PSQuery, PSConfig>;
export const PSQueryEditor = (props: PSQueryEditorProps) => {
  const { datasource, query, onChange, onRunQuery } = props;
  const { getOpenAPI3Spec } = datasource;
  const [spec, setSpec] = useState<OpenAPI3Spec | null>(null);
  useEffect(() => {
    getOpenAPI3Spec().then(setSpec).catch(console.error);
  }, [getOpenAPI3Spec]);
  const onBuilderOptionsChange = (newQuery: OpenAPI3Query, runQuery = false) => {
    onChange({ ...query, ...newQuery });
    onRunQuery();
  };
  return spec ? <OpenAPI3Editor spec={spec} query={query} onChange={onBuilderOptionsChange} /> : <>Loading...</>;
};

type PSConfigEditorProps = DataSourcePluginOptionsEditorProps<PSConfig, PSSecureConfig>;
export const PSConfigEditor = (props: PSConfigEditorProps) => <>Pet Store Config Editor</>;

type PSVariableEditorProps = { query: PSVariableQuery; onChange: (query: PSVariableQuery, definition: string) => void; datasource: PSDataSource };
export const PSVariablesEditor = (props: PSVariableEditorProps) => <>Pet Store Variable Editor</>;
