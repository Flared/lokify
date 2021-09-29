import 'styles/components/View.css'

import { Header } from 'components/Header'
import { Logs } from 'components/Logs'
import { Query } from 'components/Query'
import React, { useState } from 'react'

export const View = (): JSX.Element => {
  const [containerId] = useState<string | null>('fkff40')

  return (
    <div className="view">
      <Header containerId={containerId} />
      <Query />
      <Logs />
    </div>
  )
}
