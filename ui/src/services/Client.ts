import Axios, { AxiosResponse } from 'axios';

const client = Axios.create({
  baseURL: 'http://localhost:8080',
});

type QueryResponse = AxiosResponse<Query>;

interface Query {
  status: string;
  data: {
    resultType: string;
    result: QueryStreamValue[];
  };
}

export interface QueryStreamValue {
  stream: { [label: string]: string };
  values: Array<string[]>;
}

export async function getStatus(): Promise<string> {
  const ap: AxiosResponse<string> = await client.get('/api/status');
  return ap.data;
}

export async function getQuery(): Promise<QueryStreamValue[]> {
  const ap: QueryResponse = await client.get('/api/query', {
    params: {
      query: `{container_name="firework-api"}`,
    },
  });
  return ap.data.data.result;
}
