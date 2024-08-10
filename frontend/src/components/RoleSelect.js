import React from 'react';

const roles = ['Top', 'Jungle', 'Mid', 'ADC', 'Support'];

function RoleSelect({ onSelect }) {
  return (
    <select onChange={(e) => onSelect(e.target.value)}>
      <option value="">Select a role</option>
      {roles.map(role => (
        <option key={role} value={role.toLowerCase()}>{role}</option>
      ))}
    </select>
  );
}

export default RoleSelect;