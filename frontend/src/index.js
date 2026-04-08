import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';

// Убираем StrictMode, т.к. он может вызывать проблемы с render
const root = createRoot(document.getElementById('root'));
root.render(<App />);
