import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { ApolloClient, InMemoryCache, HttpLink, from } from '@apollo/client';
import { ApolloProvider } from '@apollo/client/react';
import { BrowserRouter } from 'react-router-dom';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import { getCurrentUser } from 'aws-amplify/auth';

// 認証ヘッダーを追加するリンク
const authLink = setContext(async (_, { headers }) => {
  try {
    // 現在のユーザーを取得
    const user = await getCurrentUser();
    
    // ユーザーが認証されている場合、トークンを取得
    if (user) {
      // Amplifyからトークンを取得
      const session = await import('aws-amplify/auth').then(module => 
        module.fetchAuthSession()
      );
      
      const token = session.tokens?.idToken?.toString();
      
      if (token) {
        return {
          headers: {
            ...headers,
            authorization: `Bearer ${token}`,
          }
        };
      }
    }
  } catch (error) {
    console.log('No authenticated user or token available:', error);
  }
  
  return {
    headers: {
      ...headers,
    }
  };
});

// エラーハンドリングリンク
const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) => {
      console.log(
        `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
      );
      
      // 認証エラーの場合、ログインページにリダイレクト
      if (message.includes('unauthorized') || message.includes('authentication')) {
        window.location.href = '/login';
      }
    });
  }

  if (networkError) {
    console.log(`[Network error]: ${networkError}`);
    
    // ネットワークエラーの場合もログインページにリダイレクト
    if (networkError.message.includes('401') || networkError.message.includes('403')) {
      window.location.href = '/login';
    }
  }
});

// HTTPリンク
const httpLink = new HttpLink({
  uri: import.meta.env.VITE_GRAPHQL_ENDPOINT || 'http://localhost:8080/query',
});

// Apollo Clientのインスタンスを作成
const client = new ApolloClient({
  link: from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      errorPolicy: 'all',
    },
    query: {
      errorPolicy: 'all',
    },
  },
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