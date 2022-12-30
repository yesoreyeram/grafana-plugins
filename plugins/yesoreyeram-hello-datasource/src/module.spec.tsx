import React from 'react';
import { render } from '@testing-library/react';
import { DataSource, Query, QueryEditor } from './module';

describe('Editors', () => {
  describe('Query Editor', () => {
    it('should render without error', () => {
      const ds = {} as DataSource;
      const query = {} as Query;
      const onChange = jest.fn();
      const onRunQuery = jest.fn();
      const result = render(<QueryEditor datasource={ds} query={query} onChange={onChange} onRunQuery={onRunQuery} />);
      expect(result.container.firstChild).toBeInTheDocument();
    });
  });
});
