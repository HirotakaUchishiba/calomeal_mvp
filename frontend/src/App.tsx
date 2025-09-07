import { Routes, Route } from 'react-router-dom';
import { SignUpPage } from './features/auth/SignUpPage';
import { LoginPage } from './features/auth/LoginPage';
import { OnboardingPage } from './features/onboarding/OnboardingPage';
import { DashboardPage } from './features/dashboard/DashboardPage';
// import { AuthProvider } from './contexts/AuthContext';
// import { RequireAuth, RequireGuest } from './components/ProtectedRoute';
// import './amplify-config'; // Amplify設定を読み込み

function App() {
  return (
    <Routes>
      {/* 開発環境用の簡易ルート */}
      <Route path="/" element={<LoginPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignUpPage />} />
      <Route path="/onboarding" element={<OnboardingPage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
    </Routes>
  );
}

export default App;