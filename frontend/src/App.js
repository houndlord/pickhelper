import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './pages/HomePage';
import MatchupsPage from './pages/MatchupsPage';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/matchups/:champion/:role" element={<MatchupsPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;