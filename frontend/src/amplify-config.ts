// frontend/src/amplify-config.ts

import { Amplify } from 'aws-amplify';

// 開発環境用の設定（実際の値は環境変数から取得）
const amplifyConfig = {
  Auth: {
    Cognito: {
      userPoolId: import.meta.env.VITE_COGNITO_USER_POOL_ID || 'ap-northeast-1_xxxxxxxxx',
      userPoolClientId: import.meta.env.VITE_COGNITO_CLIENT_ID || 'xxxxxxxxxxxxxxxxxxxxxxxxxx',
      loginWith: {
        email: true,
        username: false,
        phone: false,
      },
      signUpVerificationMethod: 'code',
      userAttributes: {
        email: {
          required: true,
        },
        name: {
          required: false,
        },
      },
    },
  },
  API: {
    GraphQL: {
      endpoint: import.meta.env.VITE_GRAPHQL_ENDPOINT || 'http://localhost:8080/query',
      region: import.meta.env.VITE_AWS_REGION || 'ap-northeast-1',
      defaultAuthMode: 'userPool',
    },
  },
};

// 本番環境用の設定
const productionConfig = {
  Auth: {
    Cognito: {
      userPoolId: import.meta.env.VITE_COGNITO_USER_POOL_ID,
      userPoolClientId: import.meta.env.VITE_COGNITO_CLIENT_ID,
      loginWith: {
        email: true,
        username: false,
        phone: false,
      },
      signUpVerificationMethod: 'code',
      userAttributes: {
        email: {
          required: true,
        },
        name: {
          required: false,
        },
      },
    },
  },
  API: {
    GraphQL: {
      endpoint: import.meta.env.VITE_GRAPHQL_ENDPOINT,
      region: import.meta.env.VITE_AWS_REGION || 'ap-northeast-1',
      defaultAuthMode: 'userPool',
    },
  },
};

// 環境に応じて設定を選択
const config = import.meta.env.PROD ? productionConfig : amplifyConfig;

// Amplifyを設定
Amplify.configure(config);

export default config;
