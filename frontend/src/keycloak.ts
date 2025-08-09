import Keycloak from 'keycloak-js';

const keycloak = new Keycloak({
  url: 'http://localhost:8080', // URL del container Keycloak
  realm: 'scheduler',           // Realm configurato
  clientId: 'scheduler',        // Client ID configurato in Keycloak
});

export default keycloak;
