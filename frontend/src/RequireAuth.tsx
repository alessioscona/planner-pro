import React from 'react';
import { useKeycloak } from '@react-keycloak/web';

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { keycloak, initialized } = useKeycloak();

  if (!initialized) return <div>Loading authentication...</div>;
  if (!keycloak.authenticated) {
    keycloak.login();
    return <div>Redirecting to login...</div>;
  }
  return <>{children}</>;
}
