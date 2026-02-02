// DevLake Azure Infrastructure
// Deploys: Resource Group resources, ACR, Key Vault, MySQL, and Container Instances

@description('Base name for all resources')
param baseName string = 'devlake'

@description('Azure region for deployment')
param location string = resourceGroup().location

@description('Unique suffix for globally unique names')
param uniqueSuffix string = uniqueString(resourceGroup().id)

@description('MySQL admin username')
param mysqlAdminUser string = 'merico'

@description('MySQL admin password')
@secure()
param mysqlAdminPassword string

@description('DevLake encryption secret (32 characters)')
@secure()
param encryptionSecret string

@description('ACR name for container images')
param acrName string = 'devlakeacr${uniqueSuffix}'

@description('Backend image name (without registry)')
param backendImage string = 'devlake-backend:latest'

@description('Config UI image name (without registry)')
param configUiImage string = 'devlake-config-ui:latest'

@description('Grafana image name (without registry)')
param grafanaImage string = 'devlake-grafana:latest'

// Container Registry
resource acr 'Microsoft.ContainerRegistry/registries@2023-07-01' = {
  name: acrName
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    adminUserEnabled: true
  }
}

// Key Vault
resource keyVault 'Microsoft.KeyVault/vaults@2023-07-01' = {
  name: '${baseName}kv${uniqueSuffix}'
  location: location
  properties: {
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: subscription().tenantId
    enableRbacAuthorization: true
  }
}

// Store secrets in Key Vault
resource dbPasswordSecret 'Microsoft.KeyVault/vaults/secrets@2023-07-01' = {
  parent: keyVault
  name: 'db-admin-password'
  properties: {
    value: mysqlAdminPassword
  }
}

resource encryptionSecretKv 'Microsoft.KeyVault/vaults/secrets@2023-07-01' = {
  parent: keyVault
  name: 'encryption-secret'
  properties: {
    value: encryptionSecret
  }
}

// MySQL Flexible Server
resource mysqlServer 'Microsoft.DBforMySQL/flexibleServers@2023-06-30' = {
  name: '${baseName}mysql${uniqueSuffix}'
  location: location
  sku: {
    name: 'Standard_B1ms'
    tier: 'Burstable'
  }
  properties: {
    version: '8.0.21'
    administratorLogin: mysqlAdminUser
    administratorLoginPassword: mysqlAdminPassword
    storage: {
      storageSizeGB: 32
    }
    backup: {
      backupRetentionDays: 7
      geoRedundantBackup: 'Disabled'
    }
  }
}

// MySQL Database
resource mysqlDatabase 'Microsoft.DBforMySQL/flexibleServers/databases@2023-06-30' = {
  parent: mysqlServer
  name: 'lake'
  properties: {
    charset: 'utf8mb4'
    collation: 'utf8mb4_unicode_ci'
  }
}

// MySQL Firewall Rule - Allow Azure Services
resource mysqlFirewallRule 'Microsoft.DBforMySQL/flexibleServers/firewallRules@2023-06-30' = {
  parent: mysqlServer
  name: 'AllowAllAzureServicesAndResourcesWithinAzureIps'
  properties: {
    startIpAddress: '0.0.0.0'
    endIpAddress: '0.0.0.0'
  }
}

// MySQL Server Configuration - Disable invisible primary key generation
// Required for DevLake migrations that drop/recreate primary keys
resource mysqlInvisiblePKConfig 'Microsoft.DBforMySQL/flexibleServers/configurations@2023-06-30' = {
  parent: mysqlServer
  name: 'sql_generate_invisible_primary_key'
  properties: {
    value: 'OFF'
    source: 'user-override'
  }
}

// Construct DB URL with required parameters
var dbUrl = 'mysql://${mysqlAdminUser}:${mysqlAdminPassword}@${mysqlServer.properties.fullyQualifiedDomainName}:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC&tls=true'

// Backend Container Instance
resource backendContainer 'Microsoft.ContainerInstance/containerGroups@2023-05-01' = {
  name: '${baseName}-backend-${uniqueSuffix}'
  location: location
  properties: {
    containers: [
      {
        name: 'devlake-backend'
        properties: {
          image: '${acr.properties.loginServer}/${backendImage}'
          ports: [
            {
              port: 8080
              protocol: 'TCP'
            }
          ]
          environmentVariables: [
            { name: 'DB_URL', secureValue: dbUrl }
            { name: 'ENCRYPTION_SECRET', secureValue: encryptionSecret }
            { name: 'PORT', value: '8080' }
            { name: 'MODE', value: 'release' }
            { name: 'PLUGIN_DIR', value: 'bin/plugins' }
            { name: 'REMOTE_PLUGIN_DIR', value: 'python/plugins' }
            { name: 'LOGGING_DIR', value: '/app/logs' }
            { name: 'TZ', value: 'UTC' }
          ]
          resources: {
            requests: {
              cpu: 2
              memoryInGB: 4
            }
          }
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: 8080
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: '${baseName}-${uniqueSuffix}'
    }
    imageRegistryCredentials: [
      {
        server: acr.properties.loginServer
        username: acr.listCredentials().username
        password: acr.listCredentials().passwords[0].value
      }
    ]
  }
  dependsOn: [
    mysqlDatabase
    mysqlFirewallRule
  ]
}

// Grafana Container Instance
resource grafanaContainer 'Microsoft.ContainerInstance/containerGroups@2023-05-01' = {
  name: '${baseName}-grafana-${uniqueSuffix}'
  location: location
  properties: {
    containers: [
      {
        name: 'devlake-grafana'
        properties: {
          image: '${acr.properties.loginServer}/${grafanaImage}'
          ports: [
            {
              port: 3000
              protocol: 'TCP'
            }
          ]
          environmentVariables: [
            { name: 'GF_SERVER_ROOT_URL', value: 'http://${baseName}-grafana-${uniqueSuffix}.${location}.azurecontainer.io:3000' }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: 2
            }
          }
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: 3000
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: '${baseName}-grafana-${uniqueSuffix}'
    }
    imageRegistryCredentials: [
      {
        server: acr.properties.loginServer
        username: acr.listCredentials().username
        password: acr.listCredentials().passwords[0].value
      }
    ]
  }
}

// Config UI Container Instance
resource configUiContainer 'Microsoft.ContainerInstance/containerGroups@2023-05-01' = {
  name: '${baseName}-ui-${uniqueSuffix}'
  location: location
  properties: {
    containers: [
      {
        name: 'devlake-config-ui'
        properties: {
          image: '${acr.properties.loginServer}/${configUiImage}'
          ports: [
            {
              port: 4000
              protocol: 'TCP'
            }
          ]
          environmentVariables: [
            // Note: Do NOT include http:// prefix - nginx.conf adds the protocol
            { name: 'DEVLAKE_ENDPOINT', value: '${backendContainer.properties.ipAddress.fqdn}:8080' }
            { name: 'GRAFANA_ENDPOINT', value: '${grafanaContainer.properties.ipAddress.fqdn}:3000' }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: 2
            }
          }
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: 4000
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: '${baseName}-ui-${uniqueSuffix}'
    }
    imageRegistryCredentials: [
      {
        server: acr.properties.loginServer
        username: acr.listCredentials().username
        password: acr.listCredentials().passwords[0].value
      }
    ]
  }
}

// Outputs
output acrLoginServer string = acr.properties.loginServer
output acrName string = acr.name
output keyVaultName string = keyVault.name
output mysqlServerName string = mysqlServer.name
output mysqlFqdn string = mysqlServer.properties.fullyQualifiedDomainName
output backendEndpoint string = 'http://${backendContainer.properties.ipAddress.fqdn}:8080'
output grafanaEndpoint string = 'http://${grafanaContainer.properties.ipAddress.fqdn}:3000'
output configUiEndpoint string = 'http://${configUiContainer.properties.ipAddress.fqdn}:4000'
