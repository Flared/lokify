import { LogItem } from 'components/library/LogItem';
import React from 'react';
import { QueryStreamValue } from 'services/Client';

interface Props {
  logs: QueryStreamValue[];
}

export const LogTable = (props: Props): JSX.Element => {
  return (
    <div>
      {props.logs.length !== 0 ? (
        props.logs.map((value: QueryStreamValue, index: number) => {
          return (
            <div key={index}>
              <LogItem logData={value} />
            </div>
          );
        })
      ) : (
        <p>No results.</p>
      )}
    </div>
  );
};
