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

export const GET_WEIGHT_LOGS_QUERY = gql`
query GetWeightLogs($date: String!) {
  weightLogs(date: $date) {
    id
    weight
    loggedAt
  }
}
`;

// Analytics Queries
export const GET_NUTRITION_SUMMARY_QUERY = gql`
  query GetNutritionSummary($date: String!) {
    nutritionSummary(date: $date) {
      date
      caloriesIntake
      caloriesBurned
      caloriesBalance
      protein
      carbohydrate
      fat
      fiber
      sugar
      sodium
      meals {
        mealType
        calories
        protein
        carbohydrate
        fat
        foodItems
      }
    }
  }
`;

export const GET_NUTRITION_TRENDS_QUERY = gql`
  query GetNutritionTrends($startDate: String!, $endDate: String!) {
    nutritionTrends(startDate: $startDate, endDate: $endDate) {
      dailySummaries {
        date
        caloriesIntake
        caloriesBurned
        caloriesBalance
        protein
        carbohydrate
        fat
        fiber
        sugar
        sodium
        meals {
          mealType
          calories
          protein
          carbohydrate
          fat
          foodItems
        }
      }
      caloriesAvg
      proteinAvg
      carbohydrateAvg
      fatAvg
      caloriesTrend
      proteinTrend
      carbohydrateTrend
      fatTrend
    }
  }
`;

export const GET_NUTRITION_INSIGHTS_QUERY = gql`
  query GetNutritionInsights($year: String!, $month: String!) {
    nutritionInsights(year: $year, month: $month) {
      year
      month
      totalCalories
      avgDailyCalories
      totalProtein
      avgDailyProtein
      totalCarbohydrate
      avgDailyCarbohydrate
      totalFat
      avgDailyFat
      topFoods
      recommendations
    }
  }
`;

export const GET_WEIGHT_PROGRESS_QUERY = gql`
  query GetWeightProgress($startDate: String!, $endDate: String!) {
    weightProgress(startDate: $startDate, endDate: $endDate) {
      startDate
      endDate
      startWeight
      endWeight
      weightChange
      weightChangePercentage
      avgWeeklyChange
      trend
      dataPoints {
        date
        weight
      }
    }
  }
`;

export const GET_CALORIE_BALANCE_QUERY = gql`
  query GetCalorieBalance($startDate: String!, $endDate: String!) {
    calorieBalance(startDate: $startDate, endDate: $endDate) {
      startDate
      endDate
      totalCaloriesIntake
      totalCaloriesBurned
      totalCalorieBalance
      avgDailyBalance
      daysInDeficit
      daysInSurplus
      deficitPercentage
      dailyBalances {
        date
        caloriesIntake
        caloriesBurned
        balance
      }
    }
  }
`;