import React from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../../services/auth';
import './Home.css';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const isAuthenticated = authService.isAuthenticated();

  return (
    <div className="home-container">
      <div className="home-content">
        <h1>Welcome to MuSync</h1>
        <p className="tagline">Manage and explore your music libraries</p>
        
        <div className="features">
          <div className="feature-card">
            <h3>Upload Libraries</h3>
            <p>Upload your music library XML files and keep them organized in one place.</p>
          </div>
          
          <div className="feature-card">
            <h3>View Tracks</h3>
            <p>Browse through all your tracks and see detailed information about each one.</p>
          </div>
          
          <div className="feature-card">
            <h3>Manage Collections</h3>
            <p>Organize your music into collections for easy access and management.</p>
          </div>
        </div>
        
        <div className="cta-buttons">
          {isAuthenticated ? (
            <>
              <button onClick={() => navigate('/libraries')} className="primary-button">
                My Libraries
              </button>
              <button onClick={() => navigate('/profile')} className="secondary-button">
                My Profile
              </button>
            </>
          ) : (
            <>
              <button onClick={() => navigate('/login')} className="primary-button">
                Login
              </button>
              <button onClick={() => navigate('/register')} className="secondary-button">
                Register
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default Home;