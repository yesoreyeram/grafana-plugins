import React, { useState } from 'react';
import { InlineFormLabel, Button, Input, LinkButton } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { VercelConfig, VercelSecureConfig } from './../types';

type VercelConfigEditorProps = DataSourcePluginOptionsEditorProps<VercelConfig, VercelSecureConfig>;

export const VercelConfigEditor = (props: VercelConfigEditorProps) => {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonFields } = options;
  const { secureJsonData = { apiToken: '' } } = options;
  const [apiUrl, setApiUrl] = useState(jsonData.apiUrl || 'https://api.vercel.com');
  const [apiToken, setApiToken] = useState('');
  const onOptionChange = <Key extends keyof VercelConfig, Value extends VercelConfig[Key]>(option: Key, value: Value) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, [option]: value },
    });
  };
  const onSecureOptionChange = <Key extends keyof VercelSecureConfig, Value extends VercelSecureConfig[Key]>(option: Key, value: Value, set: boolean) => {
    onOptionsChange({
      ...options,
      secureJsonData: { ...secureJsonData, [option]: value },
      secureJsonFields: { ...secureJsonFields, [option]: set },
    });
  };
  return (
    <>
      <div className="gf-form">
        <InlineFormLabel tooltip={'Vercel API URL. Defaults to https://api.vercel.com'} width={12}>
          Vercel API URL
        </InlineFormLabel>
        <Input value={apiUrl} placeholder={'https://api.vercel.com'} onChange={(e) => setApiUrl(e.currentTarget.value)} onBlur={() => onOptionChange('apiUrl', apiUrl)}></Input>
      </div>
      <div className="gf-form">
        <InlineFormLabel tooltip={'API Token. Get one from https://vercel.com/account/tokens'} width={12}>
          Vercel API Token
        </InlineFormLabel>
        {secureJsonFields?.apiToken ? (
          <>
            <Input type="text" value="Configured" disabled={true}></Input>
            <Button
              style={{ marginInlineStart: '10px' }}
              variant="secondary"
              className="reset-button"
              onClick={() => {
                setApiToken('');
                onSecureOptionChange('apiToken', apiToken, false);
              }}
            >
              Reset
            </Button>
          </>
        ) : (
          <>
            <Input
              type="password"
              placeholder="Vercel API Token"
              value={apiToken || ''}
              onChange={(e) => setApiToken(e.currentTarget.value)}
              onBlur={() => onSecureOptionChange('apiToken', apiToken, true)}
            />
            <LinkButton href="https://vercel.com/account/tokens" target={'_blank'}>
              Get one
            </LinkButton>
          </>
        )}
      </div>
    </>
  );
};
