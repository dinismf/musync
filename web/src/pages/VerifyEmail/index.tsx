import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import authService, { SetPasswordCredentials } from '../../services/auth';
import './VerifyEmail.css';

const VerifyEmail: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [code, setCode] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [verifying, setVerifying] = useState<boolean>(false);
  const [verified, setVerified] = useState<boolean>(false);

  useEffect(() => {
    // Extract verification code from URL query parameters
    const params = new URLSearchParams(location.search);
    const codeParam = params.get('code');
    
    if (codeParam) {
      setCode(codeParam);
      verifyEmail(codeParam);
    }
  }, [location]);

  const verifyEmail = async (verificationCode: string) => {
    setVerifying(true);
    setError(null);

    try {
      await authService.verifyEmail(verificationCode);
      setVerified(true);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to verify email. Please try again.');
    } finally {
      setVerifying(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    if (name === 'password') {
      setPassword(value);
    } else if (name === 'confirmPassword') {
      setConfirmPassword(value);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validate passwords match
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    // Validate password length
    if (password.length < 8) {
      setError('Password must be at least 8 characters long');
      return;
    }

    setLoading(true);

    try {
      const credentials: SetPasswordCredentials = {
        code,
        password
      };
      await authService.setPassword(credentials);
      navigate('/login', { state: { message: 'Password set successfully. You can now log in.' } });
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to set password. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (verifying) {
    return (
      <div className="verify-email-container">
        <div className="verify-email-form-container">
          <h2>Verifying Email</h2>
          <p>Please wait while we verify your email...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="verify-email-container">
      <div className="verify-email-form-container">
        <h2>Set Your Password</h2>
        {error && <div className="error-message">{error}</div>}
        {verified ? (
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                name="password"
                value={password}
                onChange={handleChange}
                required
                minLength={8}
              />
            </div>
            <div className="form-group">
              <label htmlFor="confirmPassword">Confirm Password</label>
              <input
                type="password"
                id="confirmPassword"
                name="confirmPassword"
                value={confirmPassword}
                onChange={handleChange}
                required
                minLength={8}
              />
            </div>
            <button type="submit" disabled={loading}>
              {loading ? 'Setting Password...' : 'Set Password'}
            </button>
          </form>
        ) : (
          <div className="error-message">
            <p>Invalid or expired verification code. Please check your email for a valid verification link.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default VerifyEmail;