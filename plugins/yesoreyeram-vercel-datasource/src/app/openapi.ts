import { VercelQuery, VercelQueryOpenApi3 } from './../types';
import { OpenAPI3Spec } from './../types/openapi';

export const buildQuery = (baseSpec: OpenAPI3Spec, newQuery: VercelQueryOpenApi3): VercelQuery => {
  const currentSpec = newQuery?.openapi3 || {};
  let serverUrl = currentSpec.servers_url;
  if (!serverUrl) {
    serverUrl = baseSpec.servers && baseSpec.servers[0].url;
  }
  const baseServer = (baseSpec.servers || []).find((s) => s.url === serverUrl);
  let u = `${baseServer?.url}${currentSpec.path}`;
  Object.keys(baseServer?.variables || {}).forEach((k) => {
    const selectedVariable = baseServer?.variables && baseServer.variables[k];
    let selectedValue = currentSpec.servers_variables && currentSpec.servers_variables[k] ? currentSpec.servers_variables[k] : selectedVariable?.default;
    u = u.replace(`{${k}}`, selectedValue || '');
  });
  let method = (currentSpec.path_method || 'GET') as typeof newQuery.method;
  return { ...newQuery, url: u, method };
};
