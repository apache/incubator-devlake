import { BrowserRouter, Routes, Route } from 'react-router-dom';

import * as Layout from './layouts';
import * as Page from './pages';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout.Base />}>
          <Route index element={<Page.Connections />} />
          <Route path="connections" element={<Page.Connections />} />
          <Route path="connection/:type" element={<Page.Connection />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
