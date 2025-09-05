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