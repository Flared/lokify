import React, { useState } from 'react'
import { Button, TextField } from '@material-ui/core'
import 'styles/components/QueryBar.scss'

interface QueryBarProps {
    query: string
}

export default function QueryBar(props: QueryBarProps): JSX.Element {
  const [query, setQuery] = useState('')
  function sendQuery() {
    setQuery('allo')
  }
  return (
    <form className="QueryBar">
      <TextField 
        className="logql-input"
        label="query" 
        defaultValue={props.query}
        multiline
      />
      <Button onClick={sendQuery}>Query</Button>
    </form>
  )
}
