import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import libraryService, { Library } from '../../services/library';
import authService from '../../services/auth';
import './Library.css';

const LibraryPage: React.FC = () => {
  const navigate = useNavigate();
  const [libraries, setLibraries] = useState<Library[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [uploadName, setUploadName] = useState<string>('');
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState<boolean>(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState<boolean>(false);
  const [deleteConfirmation, setDeleteConfirmation] = useState<number | null>(null);
  const [deleting, setDeleting] = useState<boolean>(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [deleteSuccess, setDeleteSuccess] = useState<boolean>(false);

  useEffect(() => {
    // Check if user is authenticated
    if (!authService.isAuthenticated()) {
      navigate('/login');
      return;
    }

    fetchLibraries();
  }, [navigate]);

  const fetchLibraries = async () => {
    try {
      setLoading(true);
      const data = await libraryService.getLibraries();
      setLibraries(data);
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load libraries');
      if (err.response?.status === 401) {
        // Unauthorized, redirect to login
        authService.logout();
        navigate('/login');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      // Check if file is XML
      if (!file.name.toLowerCase().endsWith('.xml')) {
        setUploadError('Please select an XML file');
        setUploadFile(null);
        return;
      }
      setUploadFile(file);
      setUploadError(null);
    }
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!uploadFile) {
      setUploadError('Please select a file to upload');
      return;
    }

    if (!uploadName.trim()) {
      setUploadError('Please enter a name for the library');
      return;
    }

    try {
      setUploading(true);
      setUploadError(null);
      setUploadSuccess(false);

      // Convert file to base64
      const base64Data = await libraryService.fileToBase64(uploadFile);

      // Upload library
      await libraryService.uploadLibrary(uploadName, base64Data);

      // Reset form
      setUploadName('');
      setUploadFile(null);
      setUploadSuccess(true);

      // Refresh libraries
      fetchLibraries();
    } catch (err: any) {
      setUploadError(err.response?.data?.error || 'Failed to upload library');
    } finally {
      setUploading(false);
    }
  };

  const handleViewLibrary = (id: number) => {
    navigate(`/library/${id}`);
  };

  const handleDeleteClick = (id: number) => {
    console.log('Delete clicked for library ID:', id);
    // Ensure id is a number
    setDeleteConfirmation(Number(id));
    setDeleteError(null);
  };

  const handleCancelDelete = () => {
    setDeleteConfirmation(null);
  };

  const handleConfirmDelete = async (id: number) => {
    try {
      setDeleting(true);
      setDeleteError(null);

      console.log('Confirming delete for library ID:', id);

      // Delete the library
      await libraryService.deleteLibrary(id);

      // Remove the library from the state
      setLibraries(libraries.filter(lib => lib.id !== id));
      setDeleteSuccess(true);

      // Hide success message after 3 seconds
      setTimeout(() => {
        setDeleteSuccess(false);
      }, 3000);
    } catch (err: any) {
      console.error('Error deleting library:', err);
      setDeleteError(err.response?.data?.error || 'Failed to delete library');
    } finally {
      setDeleting(false);
      setDeleteConfirmation(null);
    }
  };

  return (
    <div className="library-container">
      <div className="library-content">
        <h1>Music Libraries</h1>

        <div className="upload-section">
          <h2>Upload New Library</h2>
          {uploadError && <div className="error-message">{uploadError}</div>}
          {uploadSuccess && <div className="success-message">Library uploaded successfully!</div>}

          <form onSubmit={handleUpload}>
            <div className="form-group">
              <label htmlFor="name">Library Name</label>
              <input
                type="text"
                id="name"
                value={uploadName}
                onChange={(e) => setUploadName(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="file">XML File</label>
              <input
                type="file"
                id="file"
                accept=".xml"
                onChange={handleFileChange}
                required
              />
              {uploadFile && <p className="file-name">Selected: {uploadFile.name}</p>}
            </div>

            <button type="submit" disabled={uploading}>
              {uploading ? 'Uploading...' : 'Upload Library'}
            </button>
          </form>
        </div>

        <div className="libraries-section">
          <h2>Your Libraries</h2>

          {deleteSuccess && <div className="success-message">Library deleted successfully!</div>}
          {deleteError && <div className="error-message">{deleteError}</div>}

          {loading ? (
            <p>Loading libraries...</p>
          ) : error ? (
            <div className="error-message">{error}</div>
          ) : libraries.length === 0 ? (
            <p>You don't have any libraries yet. Upload one to get started!</p>
          ) : (
            <div className="libraries-list">
              {libraries.map((library) => (
                <div key={library.id} className="library-card">
                  <h3>{library.name}</h3>
                  <p>Created: {new Date(library.created_at).toLocaleDateString()}</p>
                  <div className="library-actions">
                    <button onClick={() => handleViewLibrary(library.id)}>
                      View Library
                    </button>
                    <button 
                      onClick={() => {
                        console.log('Delete button clicked for library:', library);
                        handleDeleteClick(library.id);
                      }}
                      className="delete-button"
                    >
                      Delete Library
                    </button>
                  </div>

                  {deleteConfirmation !== null && deleteConfirmation === library.id && (
                    <div className="delete-confirmation">
                      <p>Are you sure you want to delete this library? This will also delete all tracks and playlists associated with it.</p>
                      <div className="confirmation-buttons">
                        <button 
                          onClick={() => {
                            console.log('Confirm delete clicked for library:', library);
                            handleConfirmDelete(library.id);
                          }}
                          disabled={deleting}
                          className="confirm-button"
                        >
                          {deleting ? 'Deleting...' : 'Yes, Delete'}
                        </button>
                        <button 
                          onClick={handleCancelDelete}
                          disabled={deleting}
                          className="cancel-button"
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="navigation-buttons">
          <button onClick={() => navigate('/')}>Home</button>
          <button onClick={() => navigate('/profile')}>Profile</button>
        </div>
      </div>
    </div>
  );
};

export default LibraryPage;
