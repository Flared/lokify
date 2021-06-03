import React from 'react'
import ReactDOM from 'react-dom'
import View from './components/View'
import './styles/index.css'

ReactDOM.render(
    <React.StrictMode>
        Lokify
        <View query="status"/>
    </React.StrictMode>,
    document.getElementById('root'),
)
