import 'styles/components/Query.css';

import Button from '@mui/material/Button';
import { InputTypes } from 'components/library/Enums';
import { InputBar } from 'components/library/InputBar';
import React from 'react';

export const Query = (): JSX.Element => {
  return (
    <div className="query">
      <div className="query-section">
        {' '}
        <InputBar data={''} type={InputTypes.text} />{' '}
        <Button variant="contained">Run Query</Button>{' '}
      </div>
      <div className="query-section">
        <div className="query-details">
          <p>From:</p>
          <InputBar data={''} type={InputTypes.dateTime} />
          <p>To:</p>
          <InputBar data={''} type={InputTypes.dateTime} />
        </div>
        <div className="query-details">
          <Button>Save query</Button>
        </div>
      </div>
    </div>
  );
};
