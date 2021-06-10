import React, { useState } from 'react'
import { Header } from 'components/Header'
import { Query } from 'components/Query'
import { Logs } from 'components/Logs'
import 'styles/components/View.css'

export const View = (): JSX.Element => {

  const [containerId,] = useState<string | null>('fkff40')

  return (
    <div className="view">
      <Header containerId={containerId}/>
      <Query/>
      <Logs/>
    </div>
  )
}
