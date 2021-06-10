import React, { useState } from 'react'
import { Header } from 'components/Header'
import { Query } from 'components/Query'
import { Logs } from 'components/Logs'
import * as Client from 'services/Client'
import { Button } from '@material-ui/core'
import 'styles/components/View.css'

export const View = (): JSX.Element => {

  const [containerId,] = useState<string | null>('fkff40')

  const sendQueryRange = (query: string) => {
    Client.getQueryRange(query).then(((value) => {
      console.log(value)
      debugger
    }))
  }

  return (
    <div className="view">
      <Header containerId={containerId}/>
      <Button onClick={() => {sendQueryRange('{container_name="firework-api"}')}}>Send</Button>
      <Query/>
      <Logs/>
    </div>
  )
}
