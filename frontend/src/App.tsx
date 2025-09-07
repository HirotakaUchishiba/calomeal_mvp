import { Routes, Route } from 'react-router-dom';
import { SignUpPage } from './features/auth/SignUpPage';
import { LoginPage } from './features/auth/LoginPage';
import { OnboardingPage } from './features/onboarding/OnboardingPage';
import { DashboardPage } from './features/dashboard/DashboardPage';
import { AuthProvider } from './contexts/AuthContext';
import { RequireAuth, RequireGuest } from './components/ProtectedRoute';
import './amplify-config'; // Amplify設定を読み込み

function App() {
  return (
    <AuthProvider>
      <Routes>
        {/* 認証不要なルート */}
        <Route 
          path="/" 
          element={
            <RequireGuest>
              <LoginPage />
            </RequireGuest>
          } 
        />
        <Route 
          path="/login" 
          element={
            <RequireGuest>
              <LoginPage />
            </RequireGuest>
          } 
        />
        <Route 
          path="/signup" 
          element={
            <RequireGuest>
              <SignUpPage />
            </RequireGuest>
          } 
        />
        
        {/* 認証が必要なルート */}
        <Route 
          path="/onboarding" 
          element={
            <RequireAuth>
              <OnboardingPage />
            </RequireAuth>
          } 
        />
        <Route 
          path="/dashboard" 
          element={
            <RequireAuth>
              <DashboardPage />
            </RequireAuth>
          } 
        />
      </Routes>
    </AuthProvider>
  );
}

export default App;