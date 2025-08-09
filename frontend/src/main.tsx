import React from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { KeycloakProvider } from './keycloak-provider';
import App from './App';

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <KeycloakProvider>
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </KeycloakProvider>
  </React.StrictMode>
);
