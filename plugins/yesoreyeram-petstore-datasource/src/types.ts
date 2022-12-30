import { DataSourceJsonData, DataQuery } from '@grafana/data';

export type PSConfig = {} & DataSourceJsonData;
export type PSSecureConfig = {};

export type KV = { key: string; value: string };
export type PSQueryType = 'openApi3';
export type PSQueryBase<T extends PSQueryType> = { queryType: T } & DataQuery;
export type OpenAPI3Query = {
  queryType: 'openApi3';
  type: 'json';
  url: string;
  method: 'GET' | 'POST';
  headers?: KV[];
  bodyType?: 'none' | 'raw' | 'form-data' | 'x-www-form-urlencoded';
  body?: string;
  bodyContentType?: string;
  bodyForm?: KV[];
  builder_options: {
    server_url?: string;
    path?: string;
    path_method?: string;
    server_variables?: Record<string, string>;
    parameters?: Record<string, any>; //`${scope}:${in}:${name}`
    path_parameters?: Record<string, string>;
    path_operation_variables?: Record<string, string>;
  };
  rootSelector?: string;
} & PSQueryBase<'openApi3'>;
export type PSQuery = OpenAPI3Query;

export type PSVariableQuery = {};

//#region Resource Calls
export type GetResourceCallBase<P extends string, Q extends Record<string, any>, R extends unknown> = {
  path: P;
  query?: Q;
  response: R;
};
export type GetResourceCallOpenAPISpec3 = GetResourceCallBase<'openapi3', {}, OpenAPI3Spec>;
export type GetResourceCall = GetResourceCallOpenAPISpec3;

//#endregion

//#region Open API 3 Spec Types
export type OpenAPI3Spec = {
  // REQUIRED. This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses. The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string.
  openapi: string;
  // REQUIRED. Provides metadata about the API. The metadata MAY be used by tooling as required.
  info: Info;
  // An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
  servers?: Server[];
  // REQUIRED. The available paths and operations for the API.
  paths: Paths;
  // An element to hold various schemas for the specification.
  components?: any;
  // A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition. To make security optional, an empty security requirement ({}) can be included in the array.
  security?: any[];
  // A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
  tags?: Tag[];
  // Additional external documentation.
  externalDocs?: ExternalDocumentation;
};
// The object provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
export type Info = {
  // REQUIRED. The title of the API.
  title: string;
  // A short description of the API. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // A URL to the Terms of Service for the API. MUST be in the format of a URL.
  termsOfService?: string;
  // The contact information for the exposed API.
  contact?: Contact;
  // The license information for the exposed API.
  license?: License;
  // REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
  version: string;
};
// Contact information for the exposed API.
export type Contact = {
  // The identifying name of the contact person/organization.
  name?: string;
  // The URL pointing to the contact information. MUST be in the format of a URL.
  url?: string;
  // The email address of the contact person/organization. MUST be in the format of an email address.
  email?: string;
};
// License information for the exposed API.
export type License = {
  // REQUIRED. The license name used for the API.
  name: string;
  // A URL to the license used for the API. MUST be in the format of a URL.
  url?: string;
};
// An object representing a Server.
export type Server = {
  // REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
  url: string;
  // An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // A map between a variable name and its value. The value is used for substitution in the server's URL template.
  variables?: Record<string, ServerVariable>;
};
// An object representing a Server Variable for server URL template substitution.
export type ServerVariable = {
  // An enumeration of string values to be used if the substitution options are from a limited set. The array SHOULD NOT be empty.
  enum?: string[];
  // REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional. If the enum is defined, the value SHOULD exist in the enum's values.
  default: string;
  // An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
  description?: string;
};
// Holds the relative paths to the individual endpoints and their operations. The path is appended to the URL from the Server Object in order to construct the full URL. The Paths MAY be empty, due to ACL constraints.
export type Paths = Record<string, PathItem>;
// Describes the operations available on a single path. A Path Item MAY be empty, due to ACL constraints. The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
export type PathItem = {
  // Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined.
  $ref?: string;
  // An optional, string summary, intended to apply to all operations in this path.
  summary?: string;
  // An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // A definition of a GET operation on this path.
  get?: Operation;
  // A definition of a POST operation on this path.
  post?: Operation;
  // An alternative server array to service all operations in this path.
  servers?: Server[];
  // A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
  parameters?: Parameter[];
};
// Describes a single operation parameter. A unique parameter is defined by a combination of a name and location.
export type Parameter = {
  name: string;
  in: 'query' | 'path' | 'header' | 'cookie';
  description?: string;
  required?: boolean;
  deprecated?: boolean;
  allowEmptyValue?: boolean;
  style?: string;
  explode?: boolean;
  allowReserved?: boolean;
  schema?: {
    type?: 'string' | 'integer' | 'object';
    format?: string;
    default?: any;
    enum?: any[];
  };
};
// Describes a single API operation on a path.
export type Operation = {
  // A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
  tags?: string[];
  // A short summary of what the operation does.
  summary?: string;
  // A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // Additional external documentation for this operation.
  externalDocs?: ExternalDocumentation;
  // Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
  operationId?: string;
  // A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
  parameters?: Parameter[];
  // The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
  requestBody?: any;
  // REQUIRED. The list of possible responses as they are returned from executing this operation.
  responses: any;
  // A map of possible out-of band callbacks related to the parent operation. The key is a unique identifier for the Callback Object. Each value in the map is a Callback Object that describes a request that may be initiated by the API provider and the expected responses.
  callbacks?: Record<string, any>;
  // Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
  deprecated?: boolean;
  // A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. To make security optional, an empty security requirement ({}) can be included in the array. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
  security?: any;
  // An alternative server array to service this operation. If an alternative server object is specified at the Path Item Object or Root level, it will be overridden by this value.
  servers?: Server[];
};
// Allows referencing an external resource for extended documentation.
export type ExternalDocumentation = {
  // A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // REQUIRED. The URL for the target documentation. Value MUST be in the format of a URL.
  url: string;
};
// Adds metadata to a single tag that is used by the Operation Object. It is not mandatory to have a Tag Object per tag defined in the Operation Object instances.
export type Tag = {
  // REQUIRED. The name of the tag.
  name: string;
  // A short description for the tag. CommonMark syntax MAY be used for rich text representation.
  description?: string;
  // Additional external documentation for this tag.
  externalDocs?: ExternalDocumentation;
};
//#endregion
