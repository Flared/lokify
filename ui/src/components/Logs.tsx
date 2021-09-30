import 'styles/components/Logs.css';

import { LogTable } from 'components/library/LogTable';
import React from 'react';

export const Logs = (): JSX.Element => {
  return (
    <div className="logs">
      <p>Logs</p>
      <LogTable
        logs={[
          {
            stream: { data_timestamp_iso: 'eeee', message: 'eeee' },
            values: [['cake'], ['cake']],
          },
          {
            stream: { data_timestamp_iso: 'eeee', message: 'eeee' },
            values: [['cake'], ['cake']],
          },
        ]}
      />
    </div>
  );
};
