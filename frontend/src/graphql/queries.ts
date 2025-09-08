import { gql } from '@apollo/client';

export const GET_DAILY_SUMMARY_QUERY = gql`
  query GetDailySummary($date: String!) {
    dailySummary(date: $date) {
      caloriesIntake
      caloriesBurned
      protein
      carbohydrate
      fat
    }
  }
`;

export const GET_FOOD_LOGS_QUERY = gql`
  query GetFoodLogs($date: String!) {
    foodLogs(date: $date) {
      id
      foodName
      quantity
      unit
      calories
      protein
      carbohydrate
      fat
      loggedAt
    }
  }
`;

export const GET_EXERCISE_LOGS_QUERY = gql`
  query GetExerciseLogs($date: String!) {
    exerciseLogs(date: $date) {
      id
      exerciseName
      durationMinutes
      caloriesBurned
      loggedAt
    }
  }
`;

export const LOG_WEIGHT_MUTATION = gql`
  mutation LogWeight($weight: Float!, $date: String!) {
    logWeight(weight: $weight, date: $date) {
      id
      weight
      loggedAt
    }
  }
`;