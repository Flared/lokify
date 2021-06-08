import React from 'react'
import 'styles/components/Header.css'

interface Props {
  containerId: string | null,
}

export const Header = (props: Props): JSX.Element => {
  return (
    <div className="header">
      <div>
        <h6>{props.containerId}</h6>
      </div>
      <div>
        <h1>Lokify</h1>
      </div>
    </div>
  )
}
