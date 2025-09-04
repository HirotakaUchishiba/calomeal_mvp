import React, { useState } from 'react';
import { gql, useLazyQuery, useMutation } from '@apollo/client';

// バックエンドのスキーマに対応するGraphQLクエリとミューテーションを定義
const SEARCH_FOOD_QUERY = gql`
  query SearchFood($query: String!) {
    searchFood(query: $query) {
      id
      name
      brand
      calories
      protein
      carbohydrate
      fat
    }
  }
`;

const LOG_FOOD_MUTATION = gql`
  mutation LogFood($input: LogFoodInput!) {
    logFood(input: $input) {
      id
    }
  }
`;

// GraphQLの型定義。手動で定義するか、後でコード生成ツールを使います。
type FoodSearchResult = {
  id: string;
  name: string;
  brand?: string;
};

type Props = {
  isOpen: boolean;
  onClose: () => void;
  logDate: string; // 記録対象の日付をPropsで受け取る
};

export const FoodLogModal = ({ isOpen, onClose, logDate }: Props) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedFood, setSelectedFood] = useState<FoodSearchResult | null>(null);
  const [quantity, setQuantity] = useState('');
  const [unit, setUnit] = useState('g');

  // useLazyQueryフックで、任意のタイミングで実行できるクエリ関数を取得
  const [searchFood, { data: searchData, loading: searchLoading, error: searchError }] = useLazyQuery(
    SEARCH_FOOD_QUERY
  );

  // useMutationフックで、食事記録用のミューテーション関数を取得
  const [logFood, { loading: logLoading, error: logError }] = useMutation(LOG_FOOD_MUTATION);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim() === '') return;
    searchFood({ variables: { query: searchQuery } });
  };

  const handleSelectFood = (food: FoodSearchResult) => {
    setSelectedFood(food);
  };

  const handleLogFood = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFood) return;

    logFood({
      variables: {
        input: {
          foodId: selectedFood.id,
          quantity: parseFloat(quantity),
          unit: unit,
          date: logDate,
        },
      },
    }).then(() => {
      // 成功したらモーダルを閉じる
      onClose();
    }).catch((err: unknown) => {
      console.error("Failed to log food:", err);
    });
  };

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop">
      <div className="modal-content">
        <button onClick={onClose}>閉じる</button>
        <h2>食事を記録</h2>

        {!selectedFood? (
          <form onSubmit={handleSearch}>
            <label htmlFor="search">食品を検索:</label>
            <input
              id="search"
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="例: 鶏むね肉"
            />
            <button type="submit" disabled={searchLoading}>
              {searchLoading? '検索中...' : '検索'}
            </button>
            {searchError && <p style={{ color: 'red' }}>検索エラー: {searchError.message}</p>}
            <ul>
              {searchData?.searchFood.map((food: FoodSearchResult) => (
                <li key={food.id} onClick={() => handleSelectFood(food)} style={{ cursor: 'pointer' }}>
                  {food.name} {food.brand && `(${food.brand})`}
                </li>
              ))}
            </ul>
          </form>
        ) : (
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
            <button type="submit" disabled={logLoading}>
              {logLoading? '記録中...' : 'この食事を記録する'}
            </button>
            <button type="button" onClick={() => setSelectedFood(null)}>
              ← 検索に戻る
            </button>
            {logError && <p style={{ color: 'red' }}>記録エラー: {logError.message}</p>}
          </form>
        )}
      </div>
    </div>
  );
};