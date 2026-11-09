
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  auth_provider VARCHAR(20) NOT NULL DEFAULT 'local',
  avatar_url TEXT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_profiles (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL UNIQUE,
  weight_kg DECIMAL(6,2) NOT NULL DEFAULT 0,
  height_cm DECIMAL(6,2) NOT NULL DEFAULT 0,
  age INT NOT NULL DEFAULT 0,
  gender VARCHAR(30) NOT NULL DEFAULT '',
  activity_level VARCHAR(30) NOT NULL DEFAULT '',
  total_workouts INT NOT NULL DEFAULT 0,
  streak_days INT NOT NULL DEFAULT 0,
  total_km DECIMAL(10,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS notification_preferences (
  user_id VARCHAR(36) PRIMARY KEY,
  workout_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  hydration_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  nutrition_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  weekly_summary BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_notification_preferences_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS daily_activities (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  date DATE NOT NULL,
  steps INT NOT NULL DEFAULT 0,
  step_goal INT NOT NULL DEFAULT 10000,
  distance_km DECIMAL(8,2) NOT NULL DEFAULT 0,
  calories_burned INT NOT NULL DEFAULT 0,
  active_minutes INT NOT NULL DEFAULT 0,
  streak_days INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_activity_user_date (user_id,date),
  CONSTRAINT fk_activities_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS routines (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  title VARCHAR(150) NOT NULL,
  subtitle VARCHAR(255) NOT NULL DEFAULT '',
  scheduled_date DATE NOT NULL,
  day_name VARCHAR(30) NOT NULL DEFAULT '',
  status ENUM('pending','in_progress','completed') NOT NULL DEFAULT 'pending',
  completed_exercises INT NOT NULL DEFAULT 0,
  total_exercises INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_routines_user_date (user_id,scheduled_date),
  CONSTRAINT fk_routines_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS routine_exercises (
  id VARCHAR(36) PRIMARY KEY,
  routine_id VARCHAR(36) NOT NULL,
  exercise_order INT NOT NULL,
  name VARCHAR(150) NOT NULL,
  sets_summary VARCHAR(100) NOT NULL DEFAULT '',
  weight_kg DECIMAL(6,2) NOT NULL DEFAULT 0,
  is_bodyweight BOOLEAN NOT NULL DEFAULT FALSE,
  is_completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_exercises_routine_order (routine_id,exercise_order),
  CONSTRAINT fk_exercises_routine FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS runs (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  status ENUM('in_progress','paused','completed','cancelled') NOT NULL DEFAULT 'in_progress',
  start_time DATETIME NOT NULL,
  end_time DATETIME NULL,
  duration_seconds INT NOT NULL DEFAULT 0,
  distance_km DECIMAL(9,3) NOT NULL DEFAULT 0,
  current_speed_kmh DECIMAL(7,2) NOT NULL DEFAULT 0,
  avg_pace_min_km DECIMAL(7,2) NOT NULL DEFAULT 0,
  calories_burned INT NOT NULL DEFAULT 0,
  avg_heart_rate_bpm INT NOT NULL DEFAULT 0,
  elevation_gain_m INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_runs_user_start (user_id,start_time),
  KEY idx_runs_user_status (user_id,status),
  CONSTRAINT fk_runs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS run_route_points (
  id VARCHAR(36) PRIMARY KEY,
  run_id VARCHAR(36) NOT NULL,
  latitude DECIMAL(10,7) NOT NULL,
  longitude DECIMAL(10,7) NOT NULL,
  speed_kmh DECIMAL(7,2) NOT NULL DEFAULT 0,
  altitude_m DECIMAL(8,2) NOT NULL DEFAULT 0,
  sequence_order INT NOT NULL,
  timestamp DATETIME NOT NULL,
  UNIQUE KEY uk_route_sequence (run_id,sequence_order),
  CONSTRAINT fk_route_run FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS daily_nutrition (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  date DATE NOT NULL,
  water_intake_ml INT NOT NULL DEFAULT 0,
  water_goal_ml INT NOT NULL DEFAULT 2500,
  calories_consumed INT NOT NULL DEFAULT 0,
  calories_goal INT NOT NULL DEFAULT 2200,
  protein_g INT NOT NULL DEFAULT 0,
  protein_goal_g INT NOT NULL DEFAULT 180,
  carbs_g INT NOT NULL DEFAULT 0,
  carbs_goal_g INT NOT NULL DEFAULT 250,
  fat_g INT NOT NULL DEFAULT 0,
  fat_goal_g INT NOT NULL DEFAULT 70,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_nutrition_date (user_id,date),
  CONSTRAINT fk_nutrition_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS hydration_logs (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  date DATE NOT NULL,
  amount_ml INT NOT NULL,
  logged_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_hydration_user_date (user_id,date),
  CONSTRAINT fk_hydration_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS meal_logs (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  date DATE NOT NULL,
  meal_type ENUM('desayuno','almuerzo','merienda','cena') NOT NULL,
  time_str VARCHAR(10) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,
  calories INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_meals_user_date (user_id,date),
  CONSTRAINT fk_meals_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

 ALTER TABLE users
  ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE user_profiles
  ADD COLUMN gender VARCHAR(30) NOT NULL DEFAULT '',
  ADD COLUMN activity_level VARCHAR(30) NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS notification_preferences (
  user_id VARCHAR(36) PRIMARY KEY,
  workout_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  hydration_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  nutrition_reminders BOOLEAN NOT NULL DEFAULT FALSE,
  weekly_summary BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_notification_preferences_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;