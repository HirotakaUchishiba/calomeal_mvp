import React, { useState } from 'react';

// ダミーの食品データ型。後でGraphQLの型に置き換えます。
type FoodSearchResult = {
  id: string;
  name: string;
};

// モーダルの表示状態を管理するためのProps
type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export const FoodLogModal = ({ isOpen, onClose }: Props) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults]= useState<FoodSearchResult[]>([]);
  const [selectedFood, setSelectedFood]= useState<FoodSearchResult | null>(null);
  const [quantity, setQuantity] = useState('');
  const [unit, setUnit] = useState('g');

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: タスク5でAPI検索ロジックを実装します
    console.log('Searching for:', searchQuery);
    // ダミーの検索結果を表示
    setSearchResults([
      { id: '1', name: 'ごはん' },
      { id: '2', name: '鶏むね肉' },
    ]);
  };

  const handleSelectFood = (food: FoodSearchResult) => {
    setSelectedFood(food);
  };

  const handleLogFood = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: タスク5でAPI記録ロジックを実装します
    console.log('Logging:', { selectedFood, quantity, unit });
    onClose(); // 記録後にモーダルを閉じる
  };

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop">
      <div className="modal-content">
        <button onClick={onClose}>閉じる</button>
        <h2>食事を記録</h2>

        {!selectedFood? (
          // 検索ステップ
          <form onSubmit={handleSearch}>
            <label htmlFor="search">食品を検索:</label>
            <input
              id="search"
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="例: 鶏むね肉"
            />
            <button type="submit">検索</button>
            <ul>
              {searchResults.map((food) => (
                <li key={food.id} onClick={() => handleSelectFood(food)}>
                  {food.name}
                </li>
              ))}
            </ul>
          </form>
        ) : (
          // 数量入力ステップ
          <form onSubmit={handleLogFood}>
            <h3>{selectedFood.name}</h3>
            <div>
              <label htmlFor="quantity">数量:</label>
              <input
                id="quantity"
                type="number"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                required
              />
              <select value={unit} onChange={(e) => setUnit(e.target.value)}>
                <option value="g">g</option>
                <option value="個">個</option>
                <option value="皿">皿</option>
              </select>
            </div>
            <button type="submit">この食事を記録する</button>
            <button type="button" onClick={() => setSelectedFood(null)}>
              ← 検索に戻る
            </button>
          </form>
        )}
      </div>
    </div>
  );
};