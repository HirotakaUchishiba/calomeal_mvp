import { Routes, Route } from 'react-router-dom';
import { SignUpPage } from './features/auth/SignUpPage';
import { LoginPage } from './features/auth/LoginPage';
import { OnboardingPage } from './features/onboarding/OnboardingPage';

function App() {
  return (
    <Routes>
      {/* とりあえずログインページをルートURLに設定 */}
      <Route path="/" element={<LoginPage />} />
      <Route path="/signup" element={<SignUpPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/onboarding" element={<OnboardingPage />} />
      {/* TODO: フェーズ3でダッシュボードのルートを追加 */}
    </Routes>
  );
}

export default App;