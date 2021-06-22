import Axios, { AxiosResponse } from 'axios'

const client = Axios.create({
  baseURL: 'http://localhost:8082',
})

type QueryResponse = AxiosResponse<Query>

interface Query {
  status: string
  data: {
    resultType: string
    result: QueryStreamValue[]
  }
}

export interface QueryStreamValue {
  stream: { [label: string]: string }
  values: Array<string[]>
}

export async function getStatus(): Promise<string> {
  const ap: AxiosResponse<string> = await client.get('/api/status')
  return ap.data
}

export async function getQuery(): Promise<QueryStreamValue[]> {
  const ap: QueryResponse = await client.get('/api/query')
  return ap.data.data.result
}

export async function getQueryRange(query: string): Promise<QueryStreamValue[]> {
  const now = Date.now()
  const oneHourFromNow = now - (1 * 60 * 60 * 1000)
  const ap: QueryResponse = await client.get('/api/query_range', {
    params: {
      query,
      start: oneHourFromNow,
      end: now,
    }
  })
  return ap.data.data.result
}
