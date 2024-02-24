## Repository configuration
Repositories can be configured using environment variables.
Find below the lists of available configuration properties per provider.

### Vault
| Environment variable                | Default value        | Description
| ---                                 | ---                  | ---
| AUTH_SERVER_VAULT_ADDR              | localhost:8200       | The address of the Vault server
| AUTH_SERVER_VAULT_TOKEN             |                      | Vault token
| AUTH_SERVER_VAULT_BASIC_KEY         | secret/basic         | Basic authentication secret key prefix
| AUTH_SERVER_VAULT_AUTHORIZATION_KEY | secret/authorization | Authorization secret key prefix

### Aerospike
| Environment variable                    | Default value | Description
| ---                                     | ---           | ---
| AUTH_SERVER_AEROSPIKE_HOST              | localhost     | The Aerospike cluster seed host
| AUTH_SERVER_AEROSPIKE_PORT              | 3000          | The Aerospike cluster seed port
| AUTH_SERVER_AEROSPIKE_NAMESPACE         | test          | The name of the namespace containing auth details
| AUTH_SERVER_AEROSPIKE_SETNAME           | auth          | The name of the set containing auth details
| AUTH_SERVER_AEROSPIKE_BASIC_KEY         | basic         | The key of the record containing the basic authentication details
| AUTH_SERVER_AEROSPIKE_AUTHORIZATION_KEY | authorization | The key of the record containing the authorization details

### Local
| Environment variable          | Default value                      | Description
| ---                           | ---                                | ---
| AUTH_SERVER_LOCAL_CONFIG_PATH | config/local_repository_config.yml | The path to the file with the local repository configuration
