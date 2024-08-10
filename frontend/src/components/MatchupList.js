import React from 'react';
import './MatchupList.css';

function MatchupList({ matchups, champions }) {
  if (!matchups || !matchups.matchups || matchups.matchups.length === 0) {
    return <p>No matchups available.</p>;
  }

  return (
    <div className="matchup-container">
      {matchups.matchups.map((matchup, index) => {
        const championData = champions[matchup.Champion] || {};
        return (
          <div key={index} className="matchup-item">
            <img 
              src={championData.AvatarURL || ''}
              alt={matchup.Champion}
              className="champion-image"
            />
            <div className="matchup-details">
              <h3 className="champion-name">{matchup.Champion}</h3>
              <p className="win-rate">Win Rate: {matchup.WinRate}%</p>
              <p className="sample-size">Sample Size: {matchup.SampleSize}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
}

export default MatchupList;