# MuSync - Music Companion App

A companion app for electronic music enthusiasts that enhances and personalizes the music discovery process by integrating various music services and platforms.

## Features

- User Authentication (Email, Google, Social Media)
- Integration with Music Platforms (Discogs, Bandcamp, Spotify)
- Artist and Label Tracking
- Personalized Feed of New Music
- Cross-Platform Syncing
- Music Downloads (Soulseek, Deezer)

## Tech Stack

### Frontend
- Vue 3
- TypeScript
- TailwindCSS
- Pinia (State Management)

### Backend
- Go
- PostgreSQL
- MongoDB
- RESTful APIs

## Project Structure

```
musync/
├── frontend/          # Vue 3 frontend application
├── backend/           # Go backend services
└── docs/             # Project documentation
```

## Getting Started

### Prerequisites
- Node.js (v18 or higher)
- Go (v1.21 or higher)
- PostgreSQL
- MongoDB

### Development Setup

1. Clone the repository
2. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

3. Set up the backend:
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

## License

MIT 