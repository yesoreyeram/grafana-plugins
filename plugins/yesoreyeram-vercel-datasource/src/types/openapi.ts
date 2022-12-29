export type OA3ParameterObject = {
  name: string;
  in: 'query' | 'path' | 'header' | 'cookie';
  description?: string;
  required?: boolean;
  deprecated?: boolean;
  allowEmptyValue?: boolean;
  style?: string;
  explode?: boolean;
};
export type OpenAPI3SpecInfo = {
  title?: string;
  version?: string;
};
export type OpenAPI3SpecServerVariable = Record<
  string,
  {
    default: string;
    description?: string;
    enum?: string[];
  }
>;
export type OpenAPI3SpecServer = {
  url: string;
  variables?: OpenAPI3SpecServerVariable;
};
export type OpenAPI3SpecPathOperation = {
  parameters?: OA3ParameterObject[];
};
export type OpenAPI3SpecPaths = Record<
  string,
  {
    summary?: string;
    description?: string;
    get?: OpenAPI3SpecPathOperation;
    post?: OpenAPI3SpecPathOperation;
  }
>;
export type OpenAPI3Spec = {
  openapi: string;
  info: OpenAPI3SpecInfo;
  servers?: OpenAPI3SpecServer[];
  paths?: OpenAPI3SpecPaths;
};
