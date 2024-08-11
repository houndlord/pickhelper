import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import MatchupList from '../components/MatchupList';
import { getMatchups, getAllChampions } from '../api';
import './MatchupsPage.css';

function MatchupsPage() {
  const { champion, role } = useParams();
  const [matchups, setMatchups] = useState(null);
  const [champions, setChampions] = useState({});
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        setError(null);
        const [matchupsData, championsData] = await Promise.all([
          getMatchups(champion, role),
          getAllChampions()
        ]);
        setMatchups(matchupsData);
        
        if (Array.isArray(championsData)) {
          const championsObj = championsData.reduce((acc, champ) => {
            acc[champ.Name] = champ;
            return acc;
          }, {});
          setChampions(championsObj);
        } else {
          throw new Error('Unexpected champions data structure');
        }
      } catch (err) {
        console.error("Error fetching data:", err);
        setError(`Failed to load data: ${err.message}`);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, [champion, role]);

  const filteredMatchups = matchups && matchups.matchups ? 
    matchups.matchups.filter(m => 
      m.Champion.toLowerCase().includes(searchTerm.toLowerCase())
    ) : [];

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!matchups || !matchups.matchups) return <div className="error">No matchup data available</div>;

  return (
    <div className="matchups-page">
      <h1 className="page-title">Counterpicks for {champion} ({role})</h1>
      <div className="search-container">
        <input
          type="text"
          placeholder="Search counterpicks..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>
      <MatchupList matchups={{...matchups, matchups: filteredMatchups}} champions={champions} />
    </div>
  );
}

export default MatchupsPage;