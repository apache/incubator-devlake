import { Webhook } from './webhook';
import * as T from './typed';

interface Props {
  type: T.ConnectionEnum;
}

export const Connection = ({ type }: Props) => {
  switch (type) {
    case T.ConnectionEnum.webhook:
      return <Webhook />;
    default:
      return null;
  }
};
