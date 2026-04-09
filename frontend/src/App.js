import React, { useState, useEffect } from 'react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5050/api/v1';

// Функция для декодирования JWT токена
const decodeJWT = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
};

// Компонент Авторизации
const AuthPage = ({ onLogin }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    fullName: '',
    phoneNumber: ''
  });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const endpoint = isLogin ? '/auth/login' : '/auth/register';
      
      // Подготовка данных для отправки
      let body;
      if (isLogin) {
        body = { 
          email: formData.email, 
          password: formData.password 
        };
      } else {
        // Для регистрации fullName обязателен
        body = { 
          email: formData.email, 
          password: formData.password,
          full_name: formData.fullName.trim()
        };
        // phoneNumber необязателен
        if (formData.phoneNumber && formData.phoneNumber.trim()) {
          body.phone_number = formData.phoneNumber.trim();
        }
      }

      const response = await fetch(`${API_URL}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        
        // Обработка различных типов ошибок
        let errorMessage = 'Authentication failed';
        
        if (errorData.message) {
          errorMessage = errorData.message;
        } else if (errorData.error) {
          errorMessage = errorData.error;
        }
        
        // Переводим технические ошибки в понятные сообщения
        if (errorMessage.includes('failed to decode')) {
          errorMessage = 'Invalid data format. Please check your input.';
        } else if (errorMessage.includes('invalid email')) {
          errorMessage = 'Please enter a valid email address.';
        } else if (errorMessage.includes('password')) {
          errorMessage = 'Password is too short or invalid.';
        } else if (errorMessage.includes('already exists')) {
          errorMessage = 'This email is already registered.';
        } else if (errorMessage.includes('invalid email or password')) {
          errorMessage = 'Invalid email or password.';
        }
        
        throw new Error(errorMessage);
      }

      const data = await response.json();
      localStorage.setItem('accessToken', data.access_token);
      onLogin();
    } catch (err) {
      setError(err.message || 'An error occurred. Please try again.');
    }
  };

  return (
    <div className="auth-container">
      <form className="auth-form" onSubmit={handleSubmit}>
        <h1>{isLogin ? 'Login' : 'Register'}</h1>
        
        {error && <div className="error">{error}</div>}
        
        <input
          type="email"
          placeholder="Email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          required
        />
        
        <input
          type="password"
          placeholder="Password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          required
        />
        
        {!isLogin && (
          <>
            <input
              type="text"
              placeholder="Full Name"
              value={formData.fullName}
              onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
              required
            />
            <input
              type="tel"
              placeholder="Phone Number (optional)"
              value={formData.phoneNumber}
              onChange={(e) => setFormData({ ...formData, phoneNumber: e.target.value })}
            />
          </>
        )}
        
        <button type="submit">{isLogin ? 'Login' : 'Register'}</button>
        <button 
          type="button" 
          className="toggle-btn"
          onClick={() => setIsLogin(!isLogin)}
        >
          {isLogin ? 'Need an account? Register' : 'Have an account? Login'}
        </button>
      </form>
    </div>
  );
};

// Компонент Задач
const TasksPage = ({ token, onLogout }) => {
  const [tasks, setTasks] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    deadline: ''
  });

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    setIsLoading(true);
    setError('');
    try {
      // Получаем user_id из токена
      const decodedToken = decodeJWT(token);
      const userId = decodedToken?.user_id;
      
      // Добавляем фильтр по userId (бэкенд ожидает user_id с подчеркиванием)
      const url = userId ? `${API_URL}/tasks?user_id=${userId}` : `${API_URL}/tasks`;
      
      const response = await fetch(url, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) throw new Error('Failed to load tasks');

      const data = await response.json();
      setTasks(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const createTask = async (e) => {
    e.preventDefault();
    try {
      // Получаем user_id из токена
      const decodedToken = decodeJWT(token);
      const userId = decodedToken?.user_id;
      
      if (!userId) {
        throw new Error('User ID not found in token');
      }
      
      const taskData = {
        title: newTask.title,
        description: newTask.description,
        author_user_id: userId,
      };
      
      // Добавляем deadline только если он указан
      if (newTask.deadline) {
        // Конвертируем datetime-local формат в ISO 8601 с таймзоной
        const deadlineDate = new Date(newTask.deadline);
        taskData.deadline = deadlineDate.toISOString();
      }

      const response = await fetch(`${API_URL}/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(taskData),
      });

      if (!response.ok) throw new Error('Failed to create task');

      setNewTask({ title: '', description: '', deadline: '' });
      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  const toggleTaskCompletion = async (taskId, currentCompleted) => {
    try {
      const response = await fetch(`${API_URL}/tasks/${taskId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ completed: !currentCompleted }),
      });

      if (!response.ok) throw new Error('Failed to update task');

      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  const deleteTask = async (taskId) => {
    if (!window.confirm('Are you sure you want to delete this task?')) return;

    try {
      const response = await fetch(`${API_URL}/tasks/${taskId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) throw new Error('Failed to delete task');

      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
        <h1>My Tasks</h1>
        <div className="header-actions">
          <button onClick={() => window.location.hash = '#stats'} className="stats-btn">
            Show Stats
          </button>
          <button onClick={() => window.location.hash = '#users'} className="users-btn">
            Show Users
          </button>
          <button onClick={() => window.location.hash = '#profile'} className="stats-btn">
            Profile
          </button>
          <button onClick={onLogout} className="logout-btn">Logout</button>
        </div>
      </div>

      {error && <div className="error">{error}</div>}

      <form className="create-task-form" onSubmit={createTask}>
        <h2>Create New Task</h2>
        <input
          type="text"
          placeholder="Task Title"
          value={newTask.title}
          onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
          required
        />
        <textarea
          placeholder="Task Description"
          value={newTask.description}
          onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
          required
        />
        <input
          type="datetime-local"
          placeholder="Deadline (optional)"
          value={newTask.deadline}
          onChange={(e) => setNewTask({ ...newTask, deadline: e.target.value })}
        />
        <button type="submit">Create Task</button>
      </form>

      {isLoading ? (
        <div className="loading">Loading tasks...</div>
      ) : (
        <div className="tasks-list">
          {tasks.length === 0 ? (
            <p>No tasks yet. Create your first task!</p>
          ) : (
            tasks.map((task) => (
              <div key={task.id} className="task-card">
                <div className="task-header">
                  <h3>{task.title}</h3>
                  <span className={`status ${task.status_key}`}>{task.status_key}</span>
                </div>
                <p>{task.description}</p>
                <div className="task-footer">
                  <small>Deadline: {new Date(task.deadline).toLocaleString()}</small>
                  <div className="task-actions">
                    <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                      <input
                        type="checkbox"
                        checked={task.completed}
                        onChange={() => toggleTaskCompletion(task.id, task.completed)}
                        className="task-completion-checkbox"
                      />
                      <span>Completed</span>
                    </label>
                    <button onClick={() => deleteTask(task.id)} className="delete-btn">
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
};

// Компонент Пользователей
const UsersPage = ({ token, onLogout }) => {
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setIsLoading(true);
    setError('');
    try {
      const response = await fetch(`${API_URL}/users`, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) throw new Error('Failed to load users');

      const data = await response.json();
      setUsers(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const deleteUser = async (userId) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return;

    try {
      const response = await fetch(`${API_URL}/users/${userId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        let errorMessage = 'Failed to delete user';
        
        // Проверяем на ошибку foreign key constraint
        if (errorData.error && errorData.error.includes('foreign key constraint')) {
          errorMessage = 'Cannot delete user with existing tasks. Please delete all user tasks first.';
        } else if (errorData.message) {
          errorMessage = errorData.message;
        }
        
        throw new Error(errorMessage);
      }

      fetchUsers();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
        <h1>Users</h1>
        <div className="header-actions">
          <button onClick={() => window.location.hash = '#home'} className="tasks-btn">
            Show Tasks
          </button>
          <button onClick={() => window.location.hash = '#stats'} className="stats-btn">
            Show Stats
          </button>
          <button onClick={() => window.location.hash = '#profile'} className="stats-btn">
            Profile
          </button>
          <button onClick={onLogout} className="logout-btn">Logout</button>
        </div>
      </div>

      {error && <div className="error">{error}</div>}

      {isLoading ? (
        <div className="loading">Loading users...</div>
      ) : (
        <div className="users-list">
          {users.length === 0 ? (
            <p>No users found.</p>
          ) : (
            users.map((user) => (
              <div key={user.id} className="user-card">
                <div className="user-info">
                  <h3>{user.full_name}</h3>
                  <p>Email: {user.email}</p>
                  <p>Phone: {user.phone_number}</p>
                  <p>Role: {user.role_key}</p>
                </div>
                <button onClick={() => deleteUser(user.id)} className="delete-btn">
                  Delete User
                </button>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
};

// Компонент Статистики
const StatsPage = ({ token, onLogout }) => {
  const [stats, setStats] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    setIsLoading(true);
    setError('');
    try {
      // Получаем user_id из токена
      const decodedToken = decodeJWT(token);
      const userId = decodedToken?.user_id;
      
      // Добавляем фильтр по userId (бэкенд ожидает user_id с подчеркиванием)
      const url = userId ? `${API_URL}/statistics?user_id=${userId}` : `${API_URL}/statistics`;
      
      const response = await fetch(url, {
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

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
        <h1>Statistics</h1>
        <div className="header-actions">
          <button onClick={() => window.location.hash = '#home'} className="tasks-btn">
            Show Tasks
          </button>
          <button onClick={() => window.location.hash = '#users'} className="users-btn">
            Show Users
          </button>
          <button onClick={() => window.location.hash = '#profile'} className="stats-btn">
            Profile
          </button>
          <button onClick={onLogout} className="logout-btn">Logout</button>
        </div>
      </div>

      {error && <div className="error">{error}</div>}

      {isLoading ? (
        <div className="loading">Loading statistics...</div>
      ) : (
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
              {stats?.tasks_completed_rate ? `${stats.tasks_completed_rate.toFixed(1)}%` : '0%'}
            </div>
            <div className="stats-card-title">Completion Rate</div>
          </div>
        </div>
      )}
    </div>
  );
};

// Компонент Личного Кабинета
const ProfilePage = ({ token, onLogout }) => {
  const [stats, setStats] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (token) {
      fetchProfileStats();
    }
  }, [token]);

  const fetchProfileStats = async () => {
    setIsLoading(true);
    setError('');
    try {
      // Получаем user_id из токена
      const decodedToken = decodeJWT(token);
      const userId = decodedToken?.user_id;
      
      // Добавляем фильтр по userId (бэкенд ожидает user_id с подчеркиванием)
      const url = userId ? `${API_URL}/statistics?user_id=${userId}` : `${API_URL}/statistics`;
      
      const response = await fetch(url, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) throw new Error('Failed to load profile statistics');

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
      <div className="profile-container">
        <div className="profile-header">
          <h1>Personal Profile</h1>
          <div className="header-actions">
            <button onClick={() => window.location.hash = '#home'} className="tasks-btn">
              Show Tasks
            </button>
            <button onClick={() => window.location.hash = '#stats'} className="stats-btn">
              Show Stats
            </button>
            <button onClick={() => window.location.hash = '#users'} className="users-btn">
              Show Users
            </button>
            <button onClick={onLogout} className="logout-btn">Logout</button>
          </div>
        </div>
        <div className="loading">Loading profile...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="profile-container">
        <div className="profile-header">
          <h1>Personal Profile</h1>
          <div className="header-actions">
            <button onClick={() => window.location.hash = '#home'} className="tasks-btn">
              Show Tasks
            </button>
            <button onClick={() => window.location.hash = '#stats'} className="stats-btn">
              Show Stats
            </button>
            <button onClick={() => window.location.hash = '#users'} className="users-btn">
              Show Users
            </button>
            <button onClick={onLogout} className="logout-btn">Logout</button>
          </div>
        </div>
        <div className="error">Error loading profile: {error}</div>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <h1>Personal Profile</h1>
        <div className="header-actions">
          <button onClick={() => window.location.hash = '#home'} className="tasks-btn">
            Show Tasks
          </button>
          <button onClick={() => window.location.hash = '#stats'} className="stats-btn">
            Show Stats
          </button>
          <button onClick={() => window.location.hash = '#users'} className="users-btn">
            Show Users
          </button>
          <button onClick={onLogout} className="logout-btn">Logout</button>
        </div>
      </div>
      
      <div className="profile-stats">
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
              {stats?.tasks_completed_rate ? `${stats.tasks_completed_rate.toFixed(1)}%` : '0%'}
            </div>
            <div className="stats-card-title">Completion Rate</div>
          </div>
        </div>
      </div>
    </div>
  );
};

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [currentPage, setCurrentPage] = useState('home');
  const [token, setToken] = useState(localStorage.getItem('accessToken') || null);

  useEffect(() => {
    setIsAuthenticated(!!token);
    
    // Обработка хешей для навигации
    const handleHashChange = () => {
      const hash = window.location.hash;
      switch (hash) {
        case '#stats':
          setCurrentPage('stats');
          break;
        case '#users':
          setCurrentPage('users');
          break;
        case '#profile':
          setCurrentPage('profile');
          break;
        default:
          setCurrentPage('home');
      }
    };

    window.addEventListener('hashchange', handleHashChange);
    handleHashChange(); // Инициализация при загрузке

    return () => {
      window.removeEventListener('hashchange', handleHashChange);
    };
  }, [token]);

  const handleLogin = () => {
    setToken(localStorage.getItem('accessToken'));
    setIsAuthenticated(true);
    window.location.hash = '#home';
  };

  const handleLogout = () => {
    localStorage.removeItem('accessToken');
    setToken(null);
    setIsAuthenticated(false);
    window.location.hash = '';
  };

  if (!isAuthenticated) {
    return <AuthPage onLogin={handleLogin} />;
  }

  return (
    <div>
      {currentPage === 'home' && <TasksPage token={token} onLogout={handleLogout} />}
      {currentPage === 'stats' && <StatsPage token={token} onLogout={handleLogout} />}
      {currentPage === 'users' && <UsersPage token={token} onLogout={handleLogout} />}
      {currentPage === 'profile' && <ProfilePage token={token} onLogout={handleLogout} />}
    </div>
  );
}

export default App;
