import React, { useState } from 'react';
import { useMutation } from '@apollo/client';
import { LOG_WEIGHT_MUTATION } from '../../graphql/queries';

type Props = {
  isOpen: boolean;
  onClose: () => void;
  logDate: string;
};

export const WeightLogModal = ({ isOpen, onClose, logDate }: Props) => {
  const [weight, setWeight] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [logWeight] = useMutation(LOG_WEIGHT_MUTATION, {
    onCompleted: () => {
      setIsSubmitting(false);
      onClose();
      setWeight('');
    },
    onError: (error) => {
      console.error('体重記録エラー:', error);
      setIsSubmitting(false);
      alert('体重の記録に失敗しました。もう一度お試しください。');
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!weight || parseFloat(weight) <= 0) {
      alert('有効な体重を入力してください。');
      return;
    }

    setIsSubmitting(true);
    
    try {
      await logWeight({
        variables: {
          weight: parseFloat(weight),
          date: logDate,
        },
      });
    } catch (error) {
      console.error('体重記録エラー:', error);
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      onClose();
      setWeight('');
    }
  };

  if (!isOpen) return null;

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.5)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 1000
    }}>
      <div style={{
        backgroundColor: 'white',
        borderRadius: '12px',
        padding: '24px',
        width: '90%',
        maxWidth: '400px',
        boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)'
      }}>
        {/* ヘッダー */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '20px'
        }}>
          <h2 style={{
            margin: 0,
            fontSize: '20px',
            fontWeight: '600',
            color: '#333'
          }}>
            📏 体重記録
          </h2>
          <button
            onClick={handleClose}
            disabled={isSubmitting}
            style={{
              background: 'none',
              border: 'none',
              fontSize: '24px',
              cursor: isSubmitting ? 'not-allowed' : 'pointer',
              color: '#666',
              opacity: isSubmitting ? 0.5 : 1
            }}
          >
            ×
          </button>
        </div>

        {/* フォーム */}
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '20px' }}>
            <label style={{
              display: 'block',
              marginBottom: '8px',
              fontSize: '14px',
              fontWeight: '500',
              color: '#555'
            }}>
              体重 (kg)
            </label>
            <input
              type="number"
              step="0.1"
              min="0"
              max="300"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
              placeholder="例: 65.5"
              disabled={isSubmitting}
              style={{
                width: '100%',
                padding: '12px',
                border: '2px solid #e1e5e9',
                borderRadius: '8px',
                fontSize: '16px',
                outline: 'none',
                transition: 'border-color 0.2s',
                opacity: isSubmitting ? 0.6 : 1
              }}
              onFocus={(e) => {
                e.target.style.borderColor = '#4CAF50';
              }}
              onBlur={(e) => {
                e.target.style.borderColor = '#e1e5e9';
              }}
            />
          </div>

          {/* 記録日表示 */}
          <div style={{
            marginBottom: '20px',
            padding: '12px',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
            fontSize: '14px',
            color: '#666'
          }}>
            📅 記録日: {new Date(logDate).toLocaleDateString('ja-JP', {
              year: 'numeric',
              month: 'long',
              day: 'numeric'
            })}
          </div>

          {/* ボタン */}
          <div style={{
            display: 'flex',
            gap: '12px'
          }}>
            <button
              type="button"
              onClick={handleClose}
              disabled={isSubmitting}
              style={{
                flex: 1,
                padding: '12px',
                border: '2px solid #e1e5e9',
                borderRadius: '8px',
                backgroundColor: 'white',
                color: '#666',
                fontSize: '16px',
                fontWeight: '500',
                cursor: isSubmitting ? 'not-allowed' : 'pointer',
                opacity: isSubmitting ? 0.5 : 1,
                transition: 'all 0.2s'
              }}
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !weight}
              style={{
                flex: 1,
                padding: '12px',
                border: 'none',
                borderRadius: '8px',
                backgroundColor: isSubmitting || !weight ? '#ccc' : '#4CAF50',
                color: 'white',
                fontSize: '16px',
                fontWeight: '500',
                cursor: isSubmitting || !weight ? 'not-allowed' : 'pointer',
                transition: 'background-color 0.2s'
              }}
            >
              {isSubmitting ? '記録中...' : '記録する'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
