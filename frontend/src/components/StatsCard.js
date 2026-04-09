import React from 'react';

export const StatsCard = ({ title, value, color = 'blue' }) => {
  const colorClasses = {
    blue: 'stats-card-blue',
    green: 'stats-card-green',
    purple: 'stats-card-purple',
    red: 'stats-card-red',
  };

  return (
    <div className={`stats-card ${colorClasses[color]}`}>
      <div className="stats-card-value">{value}</div>
      <div className="stats-card-title">{title}</div>
    </div>
  );
};