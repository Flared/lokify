import React from 'react'
import QueryBar from 'components/QueryBar'
// import { getQuery, getStatus, QueryStreamValue } from 'services/Client'
import 'styles/components/View.css'
import LogTable from './LogTable'

interface ViewProps {
    query: string
}

export default function View(props: ViewProps): JSX.Element {
  // const [status, setStatus] = useState<string | null>(null)
  // useEffect(() => {
  //     if (!status) {
  //         getStatus().then((value) => {
  //             setStatus(value)
  //         })
  //     }
  // }, [status])

  // const [logs, setLogs] = useState<QueryStreamValue[]>([])
  // useEffect(() => {
  //     if (logs.length === 0) {
  //         getQuery().then((value) => {
  //             setLogs(value)
  //         })
  //     }
  // }, [logs.length])

  return (
    <div className="view">
      <div>
        {props.query}
      </div>
      <QueryBar query="test" />
      <LogTable />
      {/* {logs.map((value: QueryStreamValue, key: number) => {
                return (
                    <div key={key}>
                        {value.stream.data_timestamp_iso} -{' '}
                        {value.stream.message}
                    </div>
                )
            })} */}
    </div>
  )
}
