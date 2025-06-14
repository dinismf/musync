import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import libraryService, { Library, Track, Playlist } from '../../services/library';
import AudioPlayer, { RHAP_UI } from 'react-h5-audio-player';
import 'react-h5-audio-player/lib/styles.css';
import './LibraryView.css';

const LibraryView: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [library, setLibrary] = useState<Library | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'tracks' | 'playlists'>('tracks');
  const [currentTrack, setCurrentTrack] = useState<Track | null>(null);
  const [isPlaying, setIsPlaying] = useState<boolean>(false);
  const [selectedPlaylist, setSelectedPlaylist] = useState<number | null>(null);
  const [playlistTracks, setPlaylistTracks] = useState<Track[]>([]);
  const [loadingPlaylist, setLoadingPlaylist] = useState<boolean>(false);
  const [currentTime, setCurrentTime] = useState<number>(0);
  const [duration, setDuration] = useState<number>(0);
  const [volume, setVolume] = useState<number>(1);
  const [isLoadingTrack, setIsLoadingTrack] = useState<boolean>(false);
  const [trackError, setTrackError] = useState<string | null>(null);
  const [streamUrl, setStreamUrl] = useState<string>('');

  useEffect(() => {
    if (!id) {
      navigate('/libraries');
      return;
    }

    const libraryId = parseInt(id);
    if (isNaN(libraryId)) {
      navigate('/libraries');
      return;
    }

    const fetchLibraryData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Fetch library details
        const libraryData = await libraryService.getLibrary(libraryId);
        setLibrary(libraryData);

        // Fetch tracks
        const tracksData = await libraryService.getTracks(libraryId);
        setTracks(tracksData);

        // Fetch playlists
        const playlistsData = await libraryService.getPlaylists(libraryId);
        setPlaylists(playlistsData);
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load library data');
        if (err.response?.status === 401) {
          navigate('/login');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchLibraryData();
  }, [id, navigate]);

  const handlePlayTrack = async (track: Track) => {
    const isCurrentTrack = currentTrack && currentTrack.id === track.id;

    if (isCurrentTrack && isPlaying) {
      // Pause current track
      setIsPlaying(false);
    } else if (isCurrentTrack && !isPlaying) {
      // Resume current track
      setIsPlaying(true);
    } else {
      // Play new track
      try {
        // Clear any previous errors
        setTrackError(null);
        setIsLoadingTrack(true);
        setCurrentTime(0);
        setDuration(0);

        // Get the stream URL based on track's location type
        const url = await libraryService.getTrackStreamUrl(track.id);

        // Set the stream URL and track info
        setStreamUrl(url);
        setCurrentTrack(track);
        setIsPlaying(true);
      } catch (error: any) {
        console.error('Error setting up track:', error);
        setTrackError(error.message || 'Failed to set up track. Please try again.');
        // Reset current track if it's a new track that failed to load
        if (!currentTrack || currentTrack.id !== track.id) {
          setCurrentTrack(null);
        }
      } finally {
        setIsLoadingTrack(false);
      }
    }
  };


  const handleSelectPlaylist = async (playlistId: number) => {
    if (selectedPlaylist === playlistId) {
      // Deselect playlist
      setSelectedPlaylist(null);
      setPlaylistTracks([]);
      return;
    }

    try {
      setLoadingPlaylist(true);
      const tracks = await libraryService.getPlaylistTracks(playlistId);
      setPlaylistTracks(tracks);
      setSelectedPlaylist(playlistId);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load playlist tracks');
    } finally {
      setLoadingPlaylist(false);
    }
  };

  const formatDuration = (seconds: number): string => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  if (loading) {
    return <div className="library-view-container"><div className="loading">Loading library data...</div></div>;
  }

  if (error) {
    return <div className="library-view-container"><div className="error">{error}</div></div>;
  }

  if (!library) {
    return <div className="library-view-container"><div className="error">Library not found</div></div>;
  }

  return (
    <div className="library-view-container">
      <div className="library-view-header">
        <h1>{library.name}</h1>
        <div className="library-info">
          <p>Source: {library.source}</p>
          <p>Created: {new Date(library.created_at).toLocaleDateString()}</p>
        </div>
        <div className="tabs">
          <button 
            className={activeTab === 'tracks' ? 'active' : ''} 
            onClick={() => setActiveTab('tracks')}
          >
            Tracks ({tracks.length})
          </button>
          <button 
            className={activeTab === 'playlists' ? 'active' : ''} 
            onClick={() => setActiveTab('playlists')}
          >
            Playlists ({playlists.length})
          </button>
        </div>
      </div>

      <div className="library-view-content">
        {activeTab === 'tracks' && (
          <div className="tracks-container">
            {tracks.length === 0 ? (
              <p className="no-items">No tracks found in this library.</p>
            ) : (
              <table className="tracks-table">
                <thead>
                  <tr>
                    <th></th>
                    <th>Title</th>
                    <th>Artist</th>
                    <th>Album</th>
                    <th>Genre</th>
                    <th>Duration</th>
                    <th>Storage</th>
                  </tr>
                </thead>
                <tbody>
                  {tracks.map((track) => (
                    <tr 
                      key={track.id} 
                      className={currentTrack?.id === track.id ? 'playing' : ''}
                      onClick={() => handlePlayTrack(track)}
                    >
                      <td>
                        {currentTrack?.id === track.id && isLoadingTrack ? (
                          <span className="loading-icon">⏳</span>
                        ) : currentTrack?.id === track.id && isPlaying ? (
                          <span className="playing-icon">▶️</span>
                        ) : (
                          <span className="play-icon">▶</span>
                        )}
                      </td>
                      <td>{track.title}</td>
                      <td>{track.artist}</td>
                      <td>{track.album}</td>
                      <td>{track.genre}</td>
                      <td>{formatDuration(track.duration)}</td>
                      <td>
                        {track.storage_type === 'local' ? (
                          <span className="storage-icon local-storage" title="Stored locally">💻 Local</span>
                        ) : track.storage_type === 'cloud' ? (
                          <span className="storage-icon cloud-storage" title="Stored in cloud">☁️ Cloud</span>
                        ) : track.storage_type === 'stream_url' ? (
                          <span className="storage-icon stream-url" title="Direct stream URL">🔗 Stream</span>
                        ) : (
                          <span className="storage-icon unknown-storage" title="Unknown storage location">❓ Unknown</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}

        {activeTab === 'playlists' && (
          <div className="playlists-container">
            <div className="playlists-list">
              {playlists.length === 0 ? (
                <p className="no-items">No playlists found in this library.</p>
              ) : (
                <ul>
                  {playlists.map((playlist) => (
                    <li 
                      key={playlist.id} 
                      className={selectedPlaylist === playlist.id ? 'selected' : ''}
                      onClick={() => handleSelectPlaylist(playlist.id)}
                    >
                      {playlist.name}
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {selectedPlaylist && (
              <div className="playlist-tracks">
                <h3>Playlist Tracks</h3>
                {loadingPlaylist ? (
                  <p>Loading playlist tracks...</p>
                ) : playlistTracks.length === 0 ? (
                  <p className="no-items">No tracks in this playlist.</p>
                ) : (
                  <table className="tracks-table">
                    <thead>
                      <tr>
                        <th></th>
                        <th>Title</th>
                        <th>Artist</th>
                        <th>Album</th>
                        <th>Duration</th>
                        <th>Storage</th>
                      </tr>
                    </thead>
                    <tbody>
                      {playlistTracks.map((track) => (
                        <tr 
                          key={track.id} 
                          className={currentTrack?.id === track.id ? 'playing' : ''}
                          onClick={() => handlePlayTrack(track)}
                        >
                          <td>
                            {currentTrack?.id === track.id && isLoadingTrack ? (
                              <span className="loading-icon">⏳</span>
                            ) : currentTrack?.id === track.id && isPlaying ? (
                              <span className="playing-icon">▶️</span>
                            ) : (
                              <span className="play-icon">▶</span>
                            )}
                          </td>
                          <td>{track.title}</td>
                          <td>{track.artist}</td>
                          <td>{track.album}</td>
                          <td>{formatDuration(track.duration)}</td>
                          <td>
                            {track.storage_type === 'local' ? (
                              <span className="storage-icon local-storage" title="Stored locally">💻 Local</span>
                            ) : track.storage_type === 'cloud' ? (
                              <span className="storage-icon cloud-storage" title="Stored in cloud">☁️ Cloud</span>
                            ) : track.storage_type === 'stream_url' ? (
                              <span className="storage-icon stream-url" title="Direct stream URL">🔗 Stream</span>
                            ) : (
                              <span className="storage-icon unknown-storage" title="Unknown storage location">❓ Unknown</span>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            )}
          </div>
        )}
      </div>

      <div className="player-bar">
        {trackError && (
          <div className="track-error-message">
            <span className="error-icon">⚠️</span> {trackError}
            <button className="dismiss-error" onClick={() => setTrackError(null)}>✕</button>
          </div>
        )}
        {currentTrack && (
          <div className="player-content">
            <div className="track-info">
              <strong>{currentTrack.title}</strong> - {currentTrack.artist}
            </div>

            <AudioPlayer
              src={streamUrl}
              autoPlay={isPlaying}
              showSkipControls={false}
              showJumpControls={true}
              onPlay={() => setIsPlaying(true)}
              onPause={() => setIsPlaying(false)}
              onEnded={() => {
                setIsPlaying(false);
                setCurrentTime(0);
              }}
              onListen={(e: Event) => setCurrentTime((e.target as HTMLAudioElement).currentTime)}
              onLoadedMetaData={(e: Event) => setDuration((e.target as HTMLAudioElement).duration)}
              onError={(e: Event) => {
                console.error('Audio playback error:', e);
                setIsPlaying(false);

                const audio = e.target as HTMLAudioElement;
                const errorMessages = {
                  1: 'Playback aborted by the user.',
                  2: 'Network error occurred during playback.',
                  3: 'Audio decoding failed. The file may be corrupted.',
                  4: 'Audio format not supported by your browser.'
                };

                const errorCode = audio.error?.code || 0;
                const errorMessage = errorMessages[errorCode as keyof typeof errorMessages] || 'An error occurred during playback.';
                setTrackError(errorMessage);
              }}
              volume={volume}
              onVolumeChange={(e: Event) => setVolume((e.target as HTMLAudioElement).volume)}
              customAdditionalControls={[]}
              layout="horizontal"
              customProgressBarSection={[
                RHAP_UI.CURRENT_TIME,
                RHAP_UI.PROGRESS_BAR,
                RHAP_UI.DURATION
              ]}
            />
          </div>
        )}
      </div>

      <div className="navigation-buttons">
        <button onClick={() => navigate('/libraries')}>Back to Libraries</button>
        <button onClick={() => navigate('/')}>Home</button>
      </div>
    </div>
  );
};

export default LibraryView;
