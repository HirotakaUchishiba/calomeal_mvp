import React, { useState } from 'react';

type ActionItem = {
  id: string;
  label: string;
  icon: string;
  onClick: () => void;
  color?: string;
};

type Props = {
  actions: ActionItem[];
  mainIcon?: string;
  mainColor?: string;
};

export const FloatingActionButton = ({ 
  actions, 
  mainIcon = '+', 
  mainColor = '#4CAF50' 
}: Props) => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  const handleActionClick = (action: ActionItem) => {
    action.onClick();
    setIsOpen(false);
  };

  return (
    <div style={{
      position: 'fixed',
      bottom: '20px',
      right: '20px',
      zIndex: 1000
    }}>
      {/* アクションアイテム */}
      {isOpen && (
        <div style={{
          position: 'absolute',
          bottom: '70px',
          right: '0',
          display: 'flex',
          flexDirection: 'column',
          gap: '10px',
          marginBottom: '10px'
        }}>
          {actions.map((action, index) => (
            <div
              key={action.id}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '10px',
                animation: `fadeInUp 0.3s ease ${index * 0.1}s both`
              }}
            >
              {/* ラベル */}
              <div style={{
                backgroundColor: 'rgba(0, 0, 0, 0.8)',
                color: 'white',
                padding: '8px 12px',
                borderRadius: '20px',
                fontSize: '14px',
                fontWeight: '500',
                whiteSpace: 'nowrap',
                boxShadow: '0 2px 8px rgba(0, 0, 0, 0.2)'
              }}>
                {action.label}
              </div>
              
              {/* アクションボタン */}
              <button
                onClick={() => handleActionClick(action)}
                style={{
                  width: '50px',
                  height: '50px',
                  borderRadius: '50%',
                  border: 'none',
                  backgroundColor: action.color || '#2196F3',
                  color: 'white',
                  fontSize: '20px',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.3)',
                  transition: 'all 0.2s ease',
                  transform: 'scale(1)'
                }}
                onMouseOver={(e) => {
                  e.currentTarget.style.transform = 'scale(1.1)';
                  e.currentTarget.style.boxShadow = '0 6px 16px rgba(0, 0, 0, 0.4)';
                }}
                onMouseOut={(e) => {
                  e.currentTarget.style.transform = 'scale(1)';
                  e.currentTarget.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.3)';
                }}
              >
                {action.icon}
              </button>
            </div>
          ))}
        </div>
      )}

      {/* メインボタン */}
      <button
        onClick={toggleMenu}
        style={{
          width: '60px',
          height: '60px',
          borderRadius: '50%',
          border: 'none',
          backgroundColor: mainColor,
          color: 'white',
          fontSize: '24px',
          cursor: 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.3)',
          transition: 'all 0.3s ease',
          transform: isOpen ? 'rotate(45deg)' : 'rotate(0deg)'
        }}
        onMouseOver={(e) => {
          e.currentTarget.style.transform = isOpen ? 'rotate(45deg) scale(1.1)' : 'rotate(0deg) scale(1.1)';
          e.currentTarget.style.boxShadow = '0 6px 16px rgba(0, 0, 0, 0.4)';
        }}
        onMouseOut={(e) => {
          e.currentTarget.style.transform = isOpen ? 'rotate(45deg)' : 'rotate(0deg)';
          e.currentTarget.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.3)';
        }}
      >
        {mainIcon}
      </button>

      {/* オーバーレイ（メニューが開いている時） */}
      {isOpen && (
        <div
          onClick={() => setIsOpen(false)}
          style={{
            position: 'fixed',
            top: '0',
            left: '0',
            right: '0',
            bottom: '0',
            backgroundColor: 'rgba(0, 0, 0, 0.1)',
            zIndex: -1
          }}
        />
      )}

      {/* CSS アニメーション */}
      <style>
        {`
          @keyframes fadeInUp {
            from {
              opacity: 0;
              transform: translateY(20px);
            }
            to {
              opacity: 1;
              transform: translateY(0);
            }
          }
        `}
      </style>
    </div>
  );
};
