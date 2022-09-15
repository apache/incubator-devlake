import { Outlet, Link } from 'react-router-dom';

import * as S from './styled';

export const Base = () => {
  return (
    <S.Container>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/connections">Connections</Link>
          </li>
        </ul>
      </nav>
      <main>
        <Outlet />
      </main>
    </S.Container>
  );
};
