import React, { useState } from 'react';
import { GET_DAILY_SUMMARY_QUERY } from '../../graphql/queries';
import { gql, useMutation } from '@apollo/client';

const LOG_EXERCISE_MUTATION = gql`
  mutation LogExercise($input: LogExerciseInput!) {
    logExercise(input: $input) {
      id
    }
  }
`;

type Props = {
  isOpen: boolean;
  onClose: () => void;
  logDate: string; // 記録対象の日付
};

export const ExerciseLogModal = ({ isOpen, onClose, logDate }: Props) => {
  const [exerciseName, setExerciseName] = useState('');
  const [durationMinutes, setDurationMinutes] = useState('');
  const [caloriesBurned, setCaloriesBurned] = useState('');

  const [logExercise, { loading, error }] = useMutation(LOG_EXERCISE_MUTATION, {
    refetchQueries:[GET_DAILY_SUMMARY_QUERY],
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    logExercise({
      variables: {
        input: {
          exerciseName,
          durationMinutes: parseInt(durationMinutes),
          caloriesBurned: parseFloat(caloriesBurned),
          date: logDate,
        },
      },
    }).then(() => {
      onClose();
    }).catch(err => {
      console.error("Failed to log exercise:", err);
    });
  };

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop">
      <div className="modal-content">
        <button onClick={onClose}>閉じる</button>
        <h2>運動を記録</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label htmlFor="exerciseName">運動名:</label>
            <input
              id="exerciseName"
              type="text"
              value={exerciseName}
              onChange={(e) => setExerciseName(e.target.value)}
              required
            />
          </div>
          <div>
            <label htmlFor="duration">実施時間 (分):</label>
            <input
              id="duration"
              type="number"
              value={durationMinutes}
              onChange={(e) => setDurationMinutes(e.target.value)}
              required
            />
          </div>
          <div>
            <label htmlFor="calories">消費カロリー (kcal):</label>
            <input
              id="calories"
              type="number"
              value={caloriesBurned}
              onChange={(e) => setCaloriesBurned(e.target.value)}
              required
            />
          </div>
          <button type="submit" disabled={loading}>
            {loading? '記録中...' : 'この運動を記録する'}
          </button>
          {error && <p style={{ color: 'red' }}>エラー: {error.message}</p>}
        </form>
      </div>
    </div>
  );
};