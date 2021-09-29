import 'styles/components/Query.css'

import { Button } from 'components/library/Button'
import { ButtonColors, ButtonSizes, InputTypes } from 'components/library/Enums'
import { InputBar } from 'components/library/InputBar'
import React from 'react'

export const Query = (): JSX.Element => {
  return (
    <div className="query">
      <div className="query-section">
        <InputBar data={''} type={InputTypes.text} />
        <Button
          size={ButtonSizes.medium}
          text={'Run Query'}
          color={ButtonColors.green}
        />
      </div>
      <div className="query-section">
        <div className="query-details">
          <p>From:</p>
          <InputBar data={''} type={InputTypes.dateTime} />
          <p>To:</p>
          <InputBar data={''} type={InputTypes.dateTime} />
        </div>
        <div className="query-details">
          <Button
            size={ButtonSizes.medium}
            text={'Save Query'}
            color={ButtonColors.blue}
          />
        </div>
      </div>
    </div>
  )
}
