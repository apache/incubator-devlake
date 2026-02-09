# PowerShell Escaping for DB_URL

## The Problem

PowerShell interprets `&` as a command separator. This truncates the DB_URL:

```powershell
# This FAILS - PowerShell sees & as "run next command"
--environment-variables DB_URL="mysql://user:pass@host:3306/db?charset=utf8mb4&parseTime=True"
                                                                              ↑ truncated here
```

## Solution 1: cmd.exe Here-String (Recommended)

Wrap the entire `az` command in a here-string and execute via cmd.exe:

```powershell
$cmd = @"
az container create --name devlake-backend --resource-group myRG --image myacr.azurecr.io/devlake-backend:latest --registry-login-server myacr.azurecr.io --registry-username myuser --registry-password mypass --dns-name-label devlake-12345 --ports 8080 --cpu 2 --memory 4 --environment-variables DB_URL="mysql://merico:password@myserver.mysql.database.azure.com:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC&tls=true" ENCRYPTION_SECRET="a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" PORT=8080 MODE=release
"@
cmd /c $cmd
```

**Why it works:** The here-string preserves the literal content, and cmd.exe doesn't interpret `&` specially within quoted strings.

## Solution 2: Backtick Escaping

Escape each `&` with PowerShell's escape character (backtick):

```powershell
--environment-variables DB_URL="mysql://user:pass@host:3306/db?charset=utf8mb4`&parseTime=True`&loc=UTC`&tls=true"
```

**Pros:** No cmd.exe needed
**Cons:** Easy to miss an ampersand

## Solution 3: Environment Variables File

Create a file with environment variables and use `--environment-variables-file`:

**env.txt:**
```
DB_URL=mysql://user:pass@host:3306/db?charset=utf8mb4&parseTime=True&loc=UTC&tls=true
ENCRYPTION_SECRET=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
PORT=8080
MODE=release
```

**Command:**
```powershell
az container create ... --environment-variables-file env.txt
```

**Pros:** Clean, no escaping needed
**Cons:** Extra file to manage, secrets in plaintext file

## Solution 4: Use Bicep (Best)

Bicep templates handle special characters correctly. See [bicep/main.bicep](../bicep/main.bicep).

```powershell
az deployment group create --resource-group myRG --template-file main.bicep --parameters ...
```

**Pros:** No escaping issues, Infrastructure-as-Code, repeatable
**Cons:** More setup initially

## Verification

Always verify the DB_URL was set correctly:

```powershell
az container show --name devlake-backend --resource-group myRG --query containers[0].environmentVariables
```

Look for the full URL including `parseTime=True&loc=UTC&tls=true`.
