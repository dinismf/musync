import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../../services/auth';
import './Profile.css';

interface UserProfile {
  id: number;
  name: string;
  email: string;
}

const Profile: React.FC = () => {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await authService.getProfile();
        setProfile(data);
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load profile');
        if (err.response?.status === 401) {
          // Unauthorized, redirect to login
          authService.logout();
          navigate('/login');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [navigate]);

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="profile-container">
        <div className="profile-card">
          <p>Loading profile...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="profile-container">
        <div className="profile-card">
          <h2>Error</h2>
          <p className="error-message">{error}</p>
          <button onClick={() => navigate('/login')}>Go to Login</button>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-card">
        <h2>Profile</h2>
        {profile && (
          <div className="profile-details">
            <div className="profile-field">
              <label>Name:</label>
              <p>{profile.name}</p>
            </div>
            <div className="profile-field">
              <label>Email:</label>
              <p>{profile.email}</p>
            </div>
          </div>
        )}
        <div className="profile-actions">
          <button onClick={() => navigate('/')}>Home</button>
          <button onClick={handleLogout} className="logout-button">Logout</button>
        </div>
      </div>
    </div>
  );
};

export default Profile;