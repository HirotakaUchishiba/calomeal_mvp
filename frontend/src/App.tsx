import { Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { RequireAuth, RequireGuest } from './components/ProtectedRoute';
import { LoginPage } from './features/auth/LoginPage';
import { SignUpPage } from './features/auth/SignUpPage';
import { OnboardingPage } from './features/onboarding/OnboardingPage';
import { DashboardPage } from './features/dashboard/DashboardPage';
import AnalyticsPage from './features/analytics/AnalyticsPage';

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/" element={
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        } />
        <Route path="/login" element={
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        } />
        <Route path="/signup" element={
          <RequireGuest>
            <SignUpPage />
          </RequireGuest>
        } />
        <Route path="/onboarding" element={
          <RequireAuth>
            <OnboardingPage />
          </RequireAuth>
        } />
        <Route path="/dashboard" element={
          <RequireAuth>
            <DashboardPage />
          </RequireAuth>
        } />
        <Route path="/analytics" element={
          <RequireAuth>
            <AnalyticsPage />
          </RequireAuth>
        } />
      </Routes>
    </AuthProvider>
  );
}

export default App;