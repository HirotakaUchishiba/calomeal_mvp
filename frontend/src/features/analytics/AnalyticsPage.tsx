import React, { useState } from 'react';
import NutritionSummaryCard from './NutritionSummaryCard';
import NutritionTrendsChart from './NutritionTrendsChart';
import WeightProgressChart from './WeightProgressChart';
import CalorieBalanceChart from './CalorieBalanceChart';

const AnalyticsPage: React.FC = () => {
  const [selectedDate, setSelectedDate] = useState(() => {
    const today = new Date();
    return today.toISOString().split('T')[0];
  });

  const [startDate, setStartDate] = useState(() => {
    const today = new Date();
    const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
    return weekAgo.toISOString().split('T')[0];
  });

  const [endDate, setEndDate] = useState(() => {
    const today = new Date();
    return today.toISOString().split('T')[0];
  });

  const [activeTab, setActiveTab] = useState<'summary' | 'trends' | 'weight' | 'balance'>('summary');

  const tabs = [
    { id: 'summary', label: '日次サマリー', icon: '📊' },
    { id: 'trends', label: '栄養トレンド', icon: '📈' },
    { id: 'weight', label: '体重進捗', icon: '⚖️' },
    { id: 'balance', label: 'カロリーバランス', icon: '⚖️' },
  ] as const;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* ヘッダー */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">アナリティクス</h1>
          <p className="mt-2 text-gray-600">あなたの健康データを分析して、より良い生活習慣を見つけましょう</p>
        </div>

        {/* タブナビゲーション */}
        <div className="mb-6">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-2 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <span className="mr-2">{tab.icon}</span>
                  {tab.label}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* 日付選択 */}
        <div className="mb-6 bg-white p-4 rounded-lg shadow">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {activeTab === 'summary' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  日付選択
                </label>
                <input
                  type="date"
                  value={selectedDate}
                  onChange={(e) => setSelectedDate(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            )}
            
            {(activeTab === 'trends' || activeTab === 'weight' || activeTab === 'balance') && (
              <>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    開始日
                  </label>
                  <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    終了日
                  </label>
                  <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </>
            )}
          </div>
        </div>

        {/* コンテンツ */}
        <div className="space-y-6">
          {activeTab === 'summary' && (
            <NutritionSummaryCard date={selectedDate} />
          )}
          
          {activeTab === 'trends' && (
            <NutritionTrendsChart startDate={startDate} endDate={endDate} />
          )}
          
          {activeTab === 'weight' && (
            <WeightProgressChart startDate={startDate} endDate={endDate} />
          )}
          
          {activeTab === 'balance' && (
            <CalorieBalanceChart startDate={startDate} endDate={endDate} />
          )}
        </div>
      </div>
    </div>
  );
};

export default AnalyticsPage;
