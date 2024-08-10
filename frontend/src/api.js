const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';

export async function getAllChampions() {
  try {
    const response = await fetch(`${API_BASE_URL}/champions`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data.champions || []; // Return the champions array or an empty array if it doesn't exist
  } catch (error) {
    console.error("Could not fetch champions:", error);
    return []; // Return empty array if fetch fails
  }
}

export async function getMatchups(champion, role) {
  try {
    const encodedChampion = encodeURIComponent(champion);
    const url = `${API_BASE_URL}/matchups/${encodedChampion}/${role}/all`;
    console.log('Fetching matchups from:', url);
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    console.log('Received matchups data:', data);
    return data;
  } catch (error) {
    console.error("Could not fetch matchups:", error);
    throw error;
  }
}