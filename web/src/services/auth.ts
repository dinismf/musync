import axios from 'axios';

// Define the base URL for API calls
const API_URL = '/api';

// Define types for authentication
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  email: string;
  username: string;
}

export interface SetPasswordCredentials {
  code: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: number;
    email: string;
    name: string;
  };
}

// Create the authentication service
class AuthService {
  // Store the JWT token in localStorage
  setToken(token: string): void {
    localStorage.setItem('token', token);
  }

  // Get the JWT token from localStorage
  getToken(): string | null {
    return localStorage.getItem('token');
  }

  // Remove the JWT token from localStorage
  removeToken(): void {
    localStorage.removeItem('token');
  }

  // Check if the user is authenticated
  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  // Login the user
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    try {
      const response = await axios.post<AuthResponse>(`${API_URL}/auth/login`, credentials);
      this.setToken(response.data.token);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Register a new user
  async register(credentials: RegisterCredentials): Promise<any> {
    try {
      const response = await axios.post(`${API_URL}/auth/signup`, credentials);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Set password during email verification
  async setPassword(credentials: SetPasswordCredentials): Promise<AuthResponse> {
    try {
      const response = await axios.post<AuthResponse>(`${API_URL}/auth/set-password`, credentials);
      this.setToken(response.data.token);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Verify email with verification code
  async verifyEmail(code: string): Promise<any> {
    try {
      const response = await axios.get(`${API_URL}/auth/verify?code=${code}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Logout the user
  logout(): void {
    this.removeToken();
  }

  // Get the authenticated user's profile
  async getProfile(): Promise<any> {
    try {
      const response = await axios.get(`${API_URL}/auth/profile`, {
        headers: {
          Authorization: `Bearer ${this.getToken()}`
        }
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  }
}

// Create and export a singleton instance
const authService = new AuthService();
export default authService;

// Set up axios interceptor to add the token to all requests
axios.interceptors.request.use(
  (config) => {
    const token = authService.getToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
