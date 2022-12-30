import React, { useState } from 'react';
import { InlineFormLabel, Input, RadioButtonGroup, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { OpenAPI3Spec, OpenAPI3Query, Operation, Parameter } from '../types';

const buildQuery = (baseSpec: OpenAPI3Spec, newQuery: OpenAPI3Query): OpenAPI3Query => {
  const currentSpec = newQuery?.builder_options || {};
  let serverUrl = currentSpec.server_url;
  if (!serverUrl) {
    serverUrl = baseSpec.servers && baseSpec.servers[0].url;
  }
  const baseServer = (baseSpec.servers || []).find((s) => s.url === serverUrl);
  let u = `${baseServer?.url || ''}${currentSpec.path || ''}`;
  Object.keys(baseServer?.variables || {}).forEach((k) => {
    const selectedVariable = baseServer?.variables && baseServer.variables[k];
    let selectedValue = currentSpec.server_variables && currentSpec.server_variables[k] ? currentSpec.server_variables[k] : selectedVariable?.default;
    u = u.replace(`{${k}}`, selectedValue || '');
  });
  Object.keys(newQuery.builder_options?.parameters || {}).forEach((k) => {
    const [_, inValue, name] = k.split(':');
    if (inValue === 'path') {
      if (newQuery?.builder_options?.parameters) {
        const paramValue = newQuery.builder_options?.parameters[k];
        u = u.replace(`{${name}}`, paramValue);
      }
    }
  });
  let method = (currentSpec.path_method || 'GET') as typeof newQuery.method;
  try {
    let a = new URL(u);
    Object.keys(newQuery.builder_options?.parameters || {}).forEach((k) => {
      const [_, inValue, name] = k.split(':');
      if (inValue === 'query') {
        if (newQuery?.builder_options?.parameters) {
          const paramValue = newQuery.builder_options?.parameters[k];
          a.searchParams.set(name, paramValue);
        }
      }
    });
    return { ...newQuery, url: a.toString(), method };
  } catch (ex) {
    console.error(ex);
  }
  return { ...newQuery, url: u, method };
};

export const OpenAPI3Editor = (props: { spec: OpenAPI3Spec; query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { query, spec } = props;
  const onChange = (newQuery: OpenAPI3Query, runQuery = false) => {
    props.onChange(buildQuery(spec, newQuery));
  };
  return <OpenSpec3Editor spec={spec} onChange={onChange} query={query} />;
};

type OA3Path = { path: string; method: string; operation: Operation; parameters?: Parameter[] };
type OA3Parameter = { path: string; method: string; parameter: Parameter; scope: 'root' | 'operation' };

const OpenSpec3Editor = (props: { spec: OpenAPI3Spec; query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { spec, query, onChange } = props;
  const allPaths: OA3Path[] = getAllPathsFromSpec(spec);
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
      <OpenSpecServersEditor spec={spec} query={query} onChange={onChange} />
      <PathEditor spec={spec} allPaths={allPaths} query={query} onChange={onChange} />
      <ParametersEditor spec={spec} allPaths={allPaths} query={query} onChange={onChange} />
      <RootSelector query={query} onChange={onChange} />
    </div>
  );
};

const OpenSpecServersEditor = (props: { spec: OpenAPI3Spec; query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { spec, query, onChange } = props;
  const [selectedServer, setSelectedServer] = useState<string>(query.builder_options?.server_url || '');
  const serverUrls = (spec.servers || []).map((s) => ({ value: s.url, label: s.url, description: s.description }));
  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
        <InlineFormLabel width={10}>Server</InlineFormLabel>
        <Select
          options={serverUrls}
          onChange={(s) => {
            setSelectedServer(s.value!);
            onChange({ ...query, builder_options: { ...query.builder_options, server_url: s.value } });
          }}
          value={selectedServer}
        />
      </div>
    </>
  );
};

const PathEditor = (props: { spec: OpenAPI3Spec; allPaths: OA3Path[]; query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { query, onChange, allPaths } = props;
  const { builder_options = {} } = query;
  const pathOptions: Array<SelectableValue<{ path: string; method: string }>> = allPaths
    .filter((p) => p.method.toUpperCase() === 'GET')
    .map((p) => {
      return {
        value: { path: p.path, method: p.method },
        label: p.method.toUpperCase() + ': ' + p.path,
        description: p.operation.summary,
      };
    })
    .sort((a, b) => (a.value.path < b.value.path ? -1 : a.value.path > b.value.path ? 1 : 0));
  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
        <InlineFormLabel width={10}>Path</InlineFormLabel>
        <Select<{ path: string; method: string }>
          options={pathOptions}
          value={pathOptions.find((p) => (p.value?.method || 'GET').toUpperCase() === (builder_options.path_method || 'GET').toUpperCase() && p.value?.path === builder_options.path)}
          onChange={(e) => {
            if (e?.value) {
              onChange({ ...query, builder_options: { ...query.builder_options, path: e.value.path, path_method: e.value.method } });
            }
          }}
        ></Select>
      </div>
    </>
  );
};

const ParametersEditor = (props: { spec: OpenAPI3Spec; allPaths: OA3Path[]; query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { query, allPaths, onChange } = props;
  const { builder_options = {} } = query;
  const currentPath = allPaths.find((p) => p.method.toUpperCase() === (builder_options.path_method || 'GET').toUpperCase() && p.path.toUpperCase() === (builder_options.path || '').toUpperCase());
  if (!currentPath) {
    return <></>;
  }
  const allParameters: OA3Parameter[] = [];
  currentPath.parameters?.forEach((p) => {
    allParameters.push({ path: currentPath.path, method: currentPath.method, parameter: p, scope: 'root' });
  });
  currentPath.operation?.parameters?.forEach((p1) => {
    allParameters.push({ path: currentPath.path, method: currentPath.method, parameter: p1, scope: 'operation' });
  });
  return (
    <>
      {allParameters.map((p) => {
        return (
          <div key={JSON.stringify(p)} style={{ display: 'flex', justifyContent: 'start', gap: '5px' }}>
            <InlineFormLabel tooltip={p.parameter.description}>{p.parameter.name}</InlineFormLabel>
            {p?.parameter?.schema?.type === 'string' && (p?.parameter?.schema?.enum || []).length > 0 && (p?.parameter?.schema?.enum || []).length < 5 ? (
              <RadioButtonGroup
                options={p?.parameter?.schema?.enum?.map((p) => ({ label: p, value: p })) || []}
                value={builder_options?.parameters ? builder_options?.parameters[p.scope + ':' + p.parameter.in + ':' + p.parameter.name || ''] : ''}
                onChange={(e) => {
                  onChange({
                    ...query,
                    builder_options: { ...builder_options, parameters: { ...builder_options.parameters, [`${p.scope}:${p.parameter.in}:${p.parameter.name}`]: e } },
                  });
                }}
              />
            ) : p?.parameter?.schema?.type === 'string' && (p?.parameter?.schema?.enum || []).length > 0 ? (
              <Select
                allowCustomValue={true}
                options={p?.parameter?.schema?.enum?.map((p) => ({ label: p, value: p })) || []}
                value={builder_options?.parameters ? builder_options?.parameters[p.scope + ':' + p.parameter.in + ':' + p.parameter.name || ''] : ''}
                onChange={(e) => {
                  onChange({
                    ...query,
                    builder_options: { ...builder_options, parameters: { ...builder_options.parameters, [`${p.scope}:${p.parameter.in}:${p.parameter.name}`]: e.value } },
                  });
                }}
              />
            ) : (
              <Input
                placeholder={p?.parameter?.schema?.default}
                value={builder_options?.parameters ? builder_options?.parameters[p.scope + ':' + p.parameter.in + ':' + p.parameter.name || ''] : ''}
                onChange={(e) => {
                  onChange({
                    ...query,
                    builder_options: { ...builder_options, parameters: { ...builder_options.parameters, [`${p.scope}:${p.parameter.in}:${p.parameter.name}`]: e.currentTarget.value } },
                  });
                }}
              />
            )}
          </div>
        );
      })}
    </>
  );
};

const RootSelector = (props: { query: OpenAPI3Query; onChange: (value: OpenAPI3Query) => void }) => {
  const { query, onChange } = props;
  const [rootSelector, setRootSelector] = useState(query.rootSelector || '');
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }}>
      <InlineFormLabel width={10}>Root Selector</InlineFormLabel>
      <Input value={rootSelector} onChange={(e) => setRootSelector(e.currentTarget.value)} onBlur={() => onChange({ ...query, rootSelector })} />
    </div>
  );
};

const getAllPathsFromSpec = (spec: OpenAPI3Spec) => {
  const allPaths: OA3Path[] = [];
  Object.keys(spec.paths).forEach((path) => {
    const currentPath = spec.paths[path];
    Object.keys(currentPath)
      .filter((method) => method.toLowerCase() === 'get' || method.toLowerCase() === 'post')
      .forEach((method) => {
        switch (method.toLowerCase()) {
          case 'get':
            if (currentPath.get) {
              allPaths.push({ path, method, operation: currentPath.get, parameters: currentPath.parameters });
            }
            break;
          case 'post':
            if (currentPath.post) {
              allPaths.push({ path, method, operation: currentPath.post, parameters: currentPath.parameters });
            }
            break;
        }
      });
  });
  return allPaths;
};
