import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import authService, { RegisterCredentials } from '../../services/auth';
import './Register.css';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const [credentials, setCredentials] = useState<RegisterCredentials>({
    email: '',
    username: ''
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [success, setSuccess] = useState<boolean>(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await authService.register(credentials);
      setSuccess(true);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to register. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-container">
      <div className="register-form-container">
        <h2>Register</h2>
        {error && <div className="error-message">{error}</div>}
        {success ? (
          <div className="success-message">
            <h3>Registration Successful!</h3>
            <p>Please check your email for a verification link to complete your registration and set your password.</p>
            <div className="register-links">
              <p>
                <Link to="/login">Back to Login</Link>
              </p>
            </div>
          </div>
        ) : (
          <>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label htmlFor="name">Name</label>
                <input
                  type="text"
                  id="username"
                  name="username"
                  value={credentials.username}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="email">Email</label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={credentials.email}
                  onChange={handleChange}
                  required
                />
              </div>
              <button type="submit" disabled={loading}>
                {loading ? 'Registering...' : 'Register'}
              </button>
            </form>
            <div className="register-links">
              <p>
                Already have an account? <Link to="/login">Login</Link>
              </p>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default Register;
