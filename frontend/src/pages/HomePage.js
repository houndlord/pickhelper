import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getAllChampions } from '../api';
import './HomePage.css';

function HomePage() {
  const [champions, setChampions] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedChampion, setSelectedChampion] = useState('');
  const [selectedRole, setSelectedRole] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchChampions() {
      const championsData = await getAllChampions();
      setChampions(championsData);
    }
    fetchChampions();
  }, []);

  const filteredChampions = champions.filter(champ =>
    champ.Name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleSearch = () => {
    if (selectedChampion && selectedRole) {
      navigate(`/matchups/${selectedChampion}/${selectedRole}`);
    }
  };

  return (
    <div className="home-page">
      <h1>pickhelper</h1>
      <div className="search-container">
        <div className="champion-search">
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onFocus={() => setShowDropdown(true)}
            onBlur={() => setTimeout(() => setShowDropdown(false), 200)}
            placeholder="Search champion..."
          />
          {showDropdown && (
            <ul className="champion-dropdown">
              {filteredChampions.map(champ => (
                <li 
                  key={champ.Name} 
                  onClick={() => {
                    setSelectedChampion(champ.Name);
                    setSearchTerm(champ.Name);
                  }}
                >
                  <img src={champ.AvatarURL} alt={champ.Name} />
                  <span>{champ.Name}</span>
                </li>
              ))}
            </ul>
          )}
        </div>
        <select 
          value={selectedRole} 
          onChange={(e) => setSelectedRole(e.target.value)}
        >
          <option value="">Select a role</option>
          <option value="top">Top</option>
          <option value="jungle">Jungle</option>
          <option value="mid">Mid</option>
          <option value="adc">ADC</option>
          <option value="support">Support</option>
        </select>
        <button onClick={handleSearch}>Find Counterpicks</button>
      </div>
    </div>
  );
}

export default HomePage;