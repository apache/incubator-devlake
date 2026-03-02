/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package token

import (
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

// tokenPrefix returns the first n characters of a token for safe logging.
func tokenPrefix(token string, n int) string {
	if len(token) <= n {
		return token
	}
	return token[:n] + "..."
}

func refreshGitHubAppInstallationToken(tp *TokenProvider) errors.Error {
	if tp == nil || tp.conn == nil {
		return errors.Default.New("missing github connection for app token refresh")
	}
	if tp.conn.AuthMethod != models.AppKey || tp.conn.InstallationID == 0 {
		return errors.Default.New("invalid github app connection for token refresh")
	}
	if tp.conn.Endpoint == "" {
		return errors.Default.New("missing github endpoint for token refresh")
	}

	oldToken := tp.conn.Token
	if tp.logger != nil {
		expiresStr := "unknown"
		if tp.conn.TokenExpiresAt != nil {
			expiresStr = tp.conn.TokenExpiresAt.Format(time.RFC3339)
		}
		tp.logger.Info(
			"Refreshing GitHub App installation token for connection %d (installation %d), old token=%s, expires_at=%s",
			tp.conn.ID, tp.conn.InstallationID,
			tokenPrefix(oldToken, 8),
			expiresStr,
		)
	}

	apiClient := newRefreshApiClient(tp.conn.Endpoint, tp.httpClient)
	installationToken, err := tp.conn.GithubAppKey.GetInstallationAccessToken(apiClient)
	if err != nil {
		if tp.logger != nil {
			tp.logger.Error(err, "Failed to refresh GitHub App installation token for connection %d", tp.conn.ID)
		}
		return err
	}

	var expiresAt *time.Time
	if !installationToken.ExpiresAt.IsZero() {
		expiresAt = &installationToken.ExpiresAt
	}
	tp.conn.UpdateToken(installationToken.Token, "", expiresAt, nil)

	if tp.logger != nil {
		tp.logger.Info(
			"Successfully refreshed GitHub App installation token for connection %d, new token=%s, new expires_at=%s",
			tp.conn.ID,
			tokenPrefix(installationToken.Token, 8),
			installationToken.ExpiresAt.Format(time.RFC3339),
		)
	}

	persistAppToken(tp.dal, tp.conn, tp.logger)
	return nil
}

func persistAppToken(d dal.Dal, conn *models.GithubConnection, logger log.Logger) {
	if d == nil || conn == nil {
		return
	}
	if err := d.UpdateColumns(conn, []dal.DalSet{
		{ColumnName: "token", Value: conn.Token},
		{ColumnName: "token_expires_at", Value: conn.TokenExpiresAt},
	}); err != nil {
		if logger != nil {
			logger.Warn(err, "Failed to persist refreshed app installation token for connection %d", conn.ID)
		}
	} else if logger != nil {
		logger.Info("Persisted refreshed app installation token for connection %d", conn.ID)
	}
}
