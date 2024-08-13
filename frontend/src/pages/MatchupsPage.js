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
        
        if (matchupsData.error) {
          throw new Error(matchupsData.error);
        }
        
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
        setError(err.message || 'An unexpected error occurred');
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
        {loading && <div className="loading">Loading...</div>}
        {!loading && error && (
          <div className="error-message">
            <h2>Oops! No data available</h2>
            <p>We couldn't find any matchup data for {champion} in the {role} role. This might be because:</p>
            <ul>
              <li>{champion} isn't commonly played as a {role}</li>
              <li>We haven't gathered enough data for this combination yet</li>
            </ul>
            <p>You could try:</p>
            <ul>
              <li>Checking a different role for {champion}</li>
              <li>Looking up a more popular champion for the {role} role</li>
            </ul>
            <p>If you think this is a mistake, please try again later or let us know!</p>
          </div>
        )}
        {!loading && !error && matchups && matchups.matchups && matchups.matchups.length > 0 && (
          <MatchupList matchups={{...matchups, matchups: filteredMatchups}} champions={champions} />
        )}
      </div>
    );
  }
  
  export default MatchupsPage;