import React, { useState, useEffect } from 'react';
import { InlineFormLabel, Input, Select } from '@grafana/ui';
import type { QueryEditorProps } from '@grafana/data/types';
import type { VercelDS } from './../datasource';
import type { OpenAPI3Spec } from './../types/openapi';
import type { VercelQuery, VercelConfig, VercelQueryOpenApi3 } from './../types';

const buildQuery = (baseSpec: OpenAPI3Spec, newQuery: VercelQueryOpenApi3): VercelQuery => {
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
  let method = (currentSpec.path_method || 'GET') as typeof newQuery.method;
  return { ...newQuery, url: u, method };
};

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
      <OpenSpecServersEditor spec={spec} query={query} onChange={onChange} />
      <PathEditor spec={spec} query={query} onChange={onChange} />
      <RootSelector query={query} onChange={onChange} />
    </div>
  );
};

const OpenSpecServersEditor = (props: { spec: OpenAPI3Spec; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
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

const PathEditor = (props: { spec: OpenAPI3Spec; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
  const { query, onChange, spec } = props;
  const { paths = {} } = spec;
  const [selectedPath, setSelectedPath] = useState(query.builder_options?.path || '');
  const [selectedMethod, setSelectedMethod] = useState(query.builder_options?.path_method || 'get');
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
              onChange({ ...query, builder_options: { ...query.builder_options, path: e.value } });
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
              onChange({ ...query, builder_options: { ...query.builder_options, path_method: e.value } });
            }
          }}
        ></Select>
      </div>
    </>
  );
};

// const OpenSpecServersVariablesEditor = (props: { variables: OpenAPI3SpecServerVariable; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
//   const { variables: spec_variables, query, onChange } = props;
//   const { openapi3 } = query;
//   const { servers_variables = {} } = openapi3 || {};
//   const keys = Object.keys(spec_variables);
//   return (
//     <>
//       {keys.map((k) => {
//         return (
//           <div style={{ display: 'flex', justifyContent: 'space-between', gap: '5px' }} key={JSON.stringify(k)}>
//             <InlineFormLabel width={10} tooltip={spec_variables[k].description || k}>
//               {k}
//             </InlineFormLabel>
//             <Input
//               placeholder={spec_variables[k].default}
//               value={servers_variables[k]}
//               onChange={(e) => onChange({ ...query, openapi3: { ...query.openapi3, servers_variables: { ...servers_variables, [k]: e.currentTarget.value || spec_variables[k].default } } })}
//             />
//           </div>
//         );
//       })}
//     </>
//   );
// };

// const OpenSpecParamVariables = (props: { spec: OpenAPI3Spec; query: VercelQueryOpenApi3; onChange: (value: VercelQuery) => void }) => {
//   const { spec, query, onChange } = props;
//   const server_params: (OpenAPI3SpecServerVariableValue | OA3ParameterObject)[] = [];
//   const selectedServer = (spec.servers || []).find((s) => s.url === query.openapi3?.servers_url);
//   if (selectedServer) {
//     Object.keys(selectedServer.variables || {}).forEach((k) => {
//       if (selectedServer.variables && selectedServer.variables[k]) {
//         server_params.push(selectedServer.variables[k]);
//       }
//     });
//   }
//   const path_params: OA3ParameterObject[] = [];
//   if (query.openapi3?.path) {
//     const selectedPath = (spec.paths || {})[query.openapi3?.path];
//     const selectedMethod = query.openapi3?.path_method || 'get';
//     let selectedParams = (spec.paths || {})[query.openapi3?.path]['get'];
//     if (query.openapi3?.path_method === 'post') {
//       selectedParams = (spec.paths || {})[query.openapi3?.path]['post'];
//     }
//     if (selectedParams && selectedParams.parameters) {
//       selectedParams.parameters.forEach((p) => {
//         path_params.push(p);
//       });
//     }
//   }
//   return <></>;
// };

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
