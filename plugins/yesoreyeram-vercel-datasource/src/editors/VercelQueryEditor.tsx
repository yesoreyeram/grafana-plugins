import React, { useState, useEffect } from 'react';
import { InlineFormLabel, Input, Select } from '@grafana/ui';
import { buildQuery } from './../app/openapi';
import type { QueryEditorProps } from '@grafana/data/types';
import type { VercelDS } from './../datasource';
import type { OpenAPI3Spec, OpenAPI3SpecPaths, OpenAPI3SpecServer, OpenAPI3SpecServerVariable } from './../types/openapi';
import type { VercelQuery, VercelConfig, VercelQueryOpenApi3 } from './../types';

type VercelQueryEditorProps = {} & QueryEditorProps<VercelDS, VercelQuery, VercelConfig>;

const isOpenAPIQuery = (query: VercelQuery): query is VercelQueryOpenApi3 => !query.queryType || query.queryType === 'openApi3';

export const VercelQueryEditor = (props: VercelQueryEditorProps) => {
  const { datasource, query, onChange, onRunQuery } = props;
  const { getOpenAPI3Spec } = datasource;
  const [spec, setSpec] = useState<OpenAPI3Spec | null>(null);
  useEffect(() => {
    getOpenAPI3Spec().then(setSpec).catch(console.error);
  }, [getOpenAPI3Spec]);
  if (query.queryType === 'raw') {
    return <></>;
  }
  return (
    <>
      {spec === null ? (
        <>Loading Query Editor</>
      ) : (
        <OpenSpec3Editor
          spec={spec}
          query={query}
          onChange={(newQuery: VercelQuery) => {
            if (isOpenAPIQuery(newQuery)) {
              onChange(buildQuery(spec, newQuery));
              onRunQuery();
            }
          }}
        />
      )}
    </>
  );
};

const OpenSpec3Editor = (props: { spec: OpenAPI3Spec; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { spec, query, onChange } = props;
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
      <OpenSpecServersEditor servers={spec?.servers || []} query={query} onChange={onChange} />
      <PathEditor paths={spec?.paths || {}} query={query} onChange={onChange} />
      <RootSelector query={query} onChange={onChange} />
    </div>
  );
};

const OpenSpecServersEditor = (props: { servers: OpenAPI3SpecServer[]; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { servers, query, onChange } = props;
  const [server, setServer] = useState<string>(query.openapi3?.servers_url || servers[0].url);
  if (servers.length < 1) {
    return <></>;
  }
  const serverUrls = servers.map((s) => ({ value: s.url, label: s.url }));
  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
        <InlineFormLabel width={10}>Server</InlineFormLabel>
        <Select
          options={serverUrls}
          onChange={(s) => {
            setServer(s.value!);
            onChange({ ...query, openapi3: { ...query.openapi3, servers_url: s.value } });
          }}
          value={server}
        />
      </div>
      {servers
        .filter((s) => s.url === server)
        .map((s) => {
          return <OpenSpecServersVariablesEditor key={JSON.stringify(s)} variables={s.variables || {}} query={query} onChange={onChange} />;
        })}
    </>
  );
};

const OpenSpecServersVariablesEditor = (props: { variables: OpenAPI3SpecServerVariable; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { variables: spec_variables, query, onChange } = props;
  const { openapi3 } = query;
  const { servers_variables = {} } = openapi3 || {};
  const keys = Object.keys(spec_variables);
  return (
    <>
      {keys.map((k) => {
        return (
          <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }} key={JSON.stringify(k)}>
            <InlineFormLabel width={10} tooltip={spec_variables[k].description || k}>
              {k}
            </InlineFormLabel>
            <Input
              placeholder={spec_variables[k].default}
              value={servers_variables[k]}
              onChange={(e) => onChange({ ...query, openapi3: { ...query.openapi3, servers_variables: { ...servers_variables, [k]: e.currentTarget.value || spec_variables[k].default } } })}
            />
          </div>
        );
      })}
    </>
  );
};

const PathEditor = (props: { paths: OpenAPI3SpecPaths; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { query, onChange, paths } = props;
  const [selectedPath, setSelectedPath] = useState(query.openapi3?.path || '');
  const [selectedMethod, setSelectedMethod] = useState(query.openapi3?.path_method || 'get');
  const pathOptions = Object.keys(paths).map((p) => ({ value: p, description: paths[p].summary || paths[p].description || p, label: p }));
  const methodOptions = Object.keys(paths[selectedPath] || { get: {} })
    .filter((p) => p.toLowerCase() === 'get' || p.toLowerCase() === 'post')
    .map((p) => ({ value: p, label: p }));
  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
        <InlineFormLabel width={10}>Path</InlineFormLabel>
        <Select
          options={pathOptions}
          value={selectedPath}
          onChange={(e) => {
            if (e?.value) {
              setSelectedPath(e.value!);
              onChange({ ...query, openapi3: { ...query.openapi3, path: e.value } });
            }
          }}
        ></Select>
        <Select
          width={20}
          options={methodOptions}
          value={selectedMethod}
          onChange={(e) => {
            if (e?.value) {
              setSelectedMethod(e.value!);
              onChange({ ...query, openapi3: { ...query.openapi3, path_method: e.value } });
            }
          }}
        ></Select>
      </div>
    </>
  );
};

const RootSelector = (props: { query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { query, onChange } = props;
  const [rootSelector, setRootSelector] = useState(query.rootSelector || '');
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
      <InlineFormLabel width={10}>Root Selector</InlineFormLabel>
      <Input value={rootSelector} onChange={(e) => setRootSelector(e.currentTarget.value)} onBlur={() => onChange({ ...query, rootSelector })} />
    </div>
  );
};
