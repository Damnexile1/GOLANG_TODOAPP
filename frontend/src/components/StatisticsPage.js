import React, { useState, useEffect } from 'react';
import './Stats.css';

const API_URL = 'http://localhost:5050/api/v1';

export const StatisticsPage = ({ token }) => {
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [stats, setStats] = useState(null);

  useEffect(() => {
    fetchStatistics();
  }, []);

  const fetchStatistics = async () => {
    setIsLoading(true);
    setError('');
    try {
      const response = await fetch(`${API_URL}/statistics`, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) throw new Error('Failed to load statistics');

      const data = await response.json();
      setStats(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="stats-loading">
        <div className="spinner"></div>
        <p>Loading statistics...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error">
        <h2>Error Loading Statistics</h2>
        <p>{error}</p>
        <button onClick={fetchStatistics}>Retry</button>
      </div>
    );
  }

  return (
    <div className="stats-container">
      <h1>Performance Statistics</h1>
      
      {/* Статистические карточки */}
      <div className="stats-grid">
        <div className="stats-card">
          <div className="stats-card-value">{stats?.tasks_created || 0}</div>
          <div className="stats-card-title">Tasks Created</div>
        </div>
        
        <div className="stats-card">
          <div className="stats-card-value">{stats?.tasks_completed || 0}</div>
          <div className="stats-card-title">Tasks Completed</div>
        </div>
        
        <div className="stats-card">
          <div className="stats-card-value">
            {stats?.tasks_completed_rate ? `${(stats.tasks_completed_rate * 100).toFixed(1)}%` : '0%'}
          </div>
          <div className="stats-card-title">Completion Rate</div>
        </div>
      </div>

      {/* Дополнительная информация */}
      <div className="additional-stats">
        <div className="stat-item">
          <span className="stat-label">Average Completion Time:</span>
          <span className="stat-value">{stats?.tasks_average_completion_time || 'N/A'}</span>
        </div>
      </div>
    </div>
  );
};