import React from 'react';

type Props = {
  summary: {
    protein: number;
    carbohydrate: number;
    fat: number;
  };
};

export const PFCProgressBars = ({ summary }: Props) => {
  // TODO: 目標値もPropsで受け取り、プログレスバーの進捗を計算する
  return (
    <div>
      <h3>PFCバランス</h3>
      <p>タンパク質(P): {summary.protein} g</p>
      <p>炭水化物(C): {summary.carbohydrate} g</p>
      <p>脂質(F): {summary.fat} g</p>
    </div>
  );
};