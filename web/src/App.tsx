import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AdminEventPage } from './pages/AdminEventPage'
import { AdminNewPage } from './pages/AdminNewPage'
import { EntryPage } from './pages/EntryPage'
import { HomePage } from './pages/HomePage'
import { ResultsPage } from './pages/ResultsPage'
import { VotePage } from './pages/VotePage'

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gradient-to-b from-violet-50 to-white">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/e/:slug" element={<EntryPage />} />
          <Route path="/e/:slug/vote" element={<VotePage />} />
          <Route path="/e/:slug/results" element={<ResultsPage />} />
          <Route path="/admin/new" element={<AdminNewPage />} />
          <Route path="/admin/:slug" element={<AdminEventPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}
