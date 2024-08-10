import React, { useState } from 'react';

function ChampionSearch({ champions, onSelect }) {
  const [searchTerm, setSearchTerm] = useState('');

  const filteredChampions = champions.filter(champ =>
    champ.Name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="champion-search">
      <input
        type="text"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        placeholder="Search champion..."
      />
      <select onChange={(e) => onSelect(e.target.value)}>
        <option value="">Select a champion</option>
        {filteredChampions.map(champ => (
          <option key={champ.Name} value={champ.Name}>
            {champ.Name}
          </option>
        ))}
      </select>
    </div>
  );
}

export default ChampionSearch;