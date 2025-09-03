import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';
import { ApolloProvider } from '@apollo/client/react';
import { BrowserRouter } from 'react-router-dom';

// Apollo Clientのインスタンスを作成
const client = new ApolloClient({
  link: new HttpLink({
    uri: 'http://localhost:8080/query',
  }),
  cache: new InMemoryCache(),
});

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <BrowserRouter> {/* この行を追加 */}
        <App />
      </BrowserRouter> {/* この行を追加 */}
    </ApolloProvider>
  </React.StrictMode>,
);