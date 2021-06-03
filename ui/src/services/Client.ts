import Axios, { AxiosResponse } from 'axios'

const client = Axios.create({
    baseURL: 'http://localhost:8080',
})

type Setter<V> = (value: V) => void

type QueryResponse = AxiosResponse<Query>

interface Query {
	status: string;
	data: QueryData;
}

interface QueryData {
	resultType: string;
	result: QueryStreamValue[];
}

export interface QueryStreamValue {
	stream: {[label: string]: string};
	values: Array<string[]>;
}

export async function getStatus(setter: Setter<string>): Promise<void> {
    const ap: AxiosResponse<string> = await client.get('/api/status')

    setter(ap.data)
}

export async function getQuery(setter: Setter<QueryStreamValue[]>): Promise<void> {
    const ap: QueryResponse = await client.get('/api/query')

    setter(ap.data.data.result)
}
