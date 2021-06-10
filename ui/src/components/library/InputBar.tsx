import React, { useState } from 'react'
import { InputTypes } from 'components/library/Enums'
import Input from '@material-ui/core/Input'
import DateFnsUtils from '@date-io/date-fns'
import { DatePicker, TimePicker, MuiPickersUtilsProvider } from '@material-ui/pickers'

interface Props {
  data: string | null,
  type: InputTypes,
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
    return (
      <></>
    )
  }
}
