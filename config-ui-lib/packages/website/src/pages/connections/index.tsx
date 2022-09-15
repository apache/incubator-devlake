import { Link } from 'react-router-dom';

import * as S from './styled';

import WebhookIcon from '@/images/icons/webhook.svg';

export const Connections = () => {
  return (
    <S.Container>
      <h1>Connections</h1>
      <h4>
        Create and manage data connections from the following data sources or Webhooks to be used in syncing data in
        your Blueprints.
      </h4>
      <div className="item">
        <h2>Webhooks</h2>
        <h4>
          You can use Webhooks to define Issues and Deployments to be used in calculating DORA metrics. Please note:
          Webhooks cannot be created or managed in Blueprints.
        </h4>
        <ul className="list">
          <li>
            <Link to="/connection/webhook">
              <img src={WebhookIcon} width={60} alt="" />
              <span>Issue/Deployment Webhook</span>
            </Link>
          </li>
        </ul>
      </div>
    </S.Container>
  );
};
