import React, { useEffect, useState } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

export const TaskStatsBarChart = ({ tasksData }) => {
  const data = {
    labels: ['Created', 'Completed'],
    datasets: [
      {
        label: 'Tasks Count',
        data: [tasksData?.tasks_created || 0, tasksData?.tasks_completed || 0],
        backgroundColor: [
          'rgba(74, 144, 226, 0.8)',
          'rgba(126, 211, 33, 0.8)',
        ],
        borderColor: [
          'rgba(74, 144, 226, 1)',
          'rgba(126, 211, 33, 1)',
        ],
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      title: {
        display: true,
        text: 'Tasks Overview',
        font: {
          size: 16,
        },
      },
      legend: {
        position: 'top',
      },
      tooltip: {
        callbacks: {
          label: function(context) {
            return `${context.dataset.label}: ${context.parsed.y}`;
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          precision: 0,
        },
      },
      x: {
        grid: {
          display: false,
        },
      },
    },
  };

  return (
    <div style={{ height: '300px', width: '100%' }}>
      <Bar data={data} options={options} />
    </div>
  );
};