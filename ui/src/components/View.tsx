import 'styles/components/View.css';

import { Header } from 'components/Header';
import { LogsTable } from 'components/LogsTable';
import { Query } from 'components/Query';
import { useEffect, useState } from 'react';
import React from 'react';
import { getQuery, QueryStreamValue } from 'services/Client';

export const View = (): JSX.Element => {
  const [containerId] = useState<string | null>('fkff40');

  const [logs, setLogs] = useState<QueryStreamValue[] | null>(null);

  useEffect(() => {
    getQuery().then((logs: QueryStreamValue[]) => {
      setLogs(logs);
    });
  });
  console.log(logs);

  return (
    <div className="view">
      <Header containerId={containerId} />
      <Query />
      <LogsTable />
    </div>
  );
};
