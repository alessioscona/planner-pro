import { useKeycloak } from '@react-keycloak/web';

export function useAuthFetch() {
  const { keycloak } = useKeycloak();

  return async (input: RequestInfo, init: RequestInit = {}) => {
    const headers = new Headers(init.headers);
    if (keycloak?.token) {
      headers.set('Authorization', `Bearer ${keycloak.token}`);
    }
    return fetch(input, { ...init, headers });
  };
}
