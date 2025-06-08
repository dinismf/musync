import axios from 'axios';

// Define the base URL for API calls
const API_URL = '/api';

// Define types for library
export interface Library {
  id: number;
  name: string;
  user_id: number;
  created_at: string;
  updated_at: string;
}

export interface Track {
  id: number;
  library_id: number;
  title: string;
  artist: string;
  album: string;
  genre: string;
  duration: number;
  year: number | null;
  created_at: string;
  updated_at: string;
}

// Create the library service
class LibraryService {
  // Get all libraries for the authenticated user
  async getLibraries(): Promise<Library[]> {
    try {
      const response = await axios.get<Library[]>(`${API_URL}/libraries`);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Get a specific library by ID
  async getLibrary(id: number): Promise<Library> {
    try {
      const response = await axios.get<Library>(`${API_URL}/libraries/${id}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Get all tracks in a library
  async getTracks(libraryId: number): Promise<Track[]> {
    try {
      const response = await axios.get<Track[]>(`${API_URL}/libraries/${libraryId}/tracks`);
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Upload a library XML file
  async uploadLibrary(name: string, fileData: string): Promise<{ library_id: number }> {
    try {
      const response = await axios.post(`${API_URL}/libraries`, {
        name,
        file_data: fileData // Base64 encoded XML file
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  // Convert a File object to base64
  async fileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        if (reader.result) {
          // Remove the data URL prefix (e.g., "data:application/xml;base64,")
          const base64String = reader.result.toString().split(',')[1];
          resolve(base64String);
        } else {
          reject(new Error('Failed to convert file to base64'));
        }
      };
      reader.onerror = error => reject(error);
    });
  }
}

// Create and export a singleton instance
const libraryService = new LibraryService();
export default libraryService;