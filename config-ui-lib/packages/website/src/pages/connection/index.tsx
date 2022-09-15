import { useParams } from 'react-router-dom';
import { Connection as ConfigConnection, ConnectionEnum } from '@devlake/config-ui';

export const Connection = () => {
  const { type } = useParams();
  const types = Object.values(ConnectionEnum) as string[];

  if (!type || !types.includes(type)) {
    return <div>someting error.</div>;
  }

  return <ConfigConnection type={type as ConnectionEnum} />;
};
