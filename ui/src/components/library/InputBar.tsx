import DateFnsUtils from '@date-io/date-fns'
import Input from '@material-ui/core/Input'
import {
  DatePicker,
  MuiPickersUtilsProvider,
  TimePicker,
} from '@material-ui/pickers'
import { InputTypes } from 'components/library/Enums'
import React, { useState } from 'react'

interface Props {
  data: string | null
  type: InputTypes
}

export const InputBar = (props: Props): JSX.Element => {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)

  if (props.type === InputTypes.text) {
    return (
      <div>
        <Input type="text" />
      </div>
    )
  } else if (props.type === InputTypes.dateTime) {
    return (
      <div>
        <MuiPickersUtilsProvider utils={DateFnsUtils}>
          <DatePicker value={selectedDate} onChange={setSelectedDate} />
          <TimePicker value={selectedDate} onChange={setSelectedDate} />
        </MuiPickersUtilsProvider>
      </div>
    )
  } else {
    return <></>
  }
}
