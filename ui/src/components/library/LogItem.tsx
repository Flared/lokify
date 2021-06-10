import React, { useState } from 'react'
import { QueryStreamValue } from 'services/Client'

interface Props {
  logData: QueryStreamValue,
}

export const LogItem = (props: Props): JSX.Element => {

  const [isExpanded, setExpanded] = useState<boolean>(false)

  const toggleLogLabelsExpand = (): void => {
    setExpanded(!isExpanded)
  }

  return (
    <div>
      <div onClick={toggleLogLabelsExpand}>
        <h6>
          {props.logData.stream.data_timestamp_iso} --- {props.logData.stream.message}
        </h6>
      </div>
      {
        isExpanded ?
        // add the log labels
          <p>Add log labels here</p>
          :
          <></>
      }
    </div>
  )
}
