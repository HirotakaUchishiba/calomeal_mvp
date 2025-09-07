// frontend/src/amplify-config.ts

import { Amplify } from 'aws-amplify';

// 開発環境用の設定（認証機能を無効化）
const amplifyConfig = {
  Auth: {
    Cognito: {
      userPoolId: 'ap-northeast-1_xxxxxxxxx',
      userPoolClientId: 'xxxxxxxxxxxxxxxxxxxxxxxxxx',
      loginWith: {
        email: true,
        username: false,
        phone: false,
      },
      signUpVerificationMethod: 'code' as const,
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
      defaultAuthMode: 'none' as const,
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
      signUpVerificationMethod: 'code' as const,
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
      defaultAuthMode: 'none' as const,
    },
  },
};

// 環境に応じて設定を選択
const config = import.meta.env.PROD ? productionConfig : amplifyConfig;

// Amplifyを設定
Amplify.configure(config);

export default config;
