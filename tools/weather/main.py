import csv
import os
import pandas as pd
import joblib
from sklearn.preprocessing import LabelEncoder
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestRegressor
from sklearn.impute import SimpleImputer

DATA_DIR = "data"
TRAINED_MODELS_DIR = "trained_models" 
WEATHER_DATA = "weather_data.csv"

def train_models(city_name: str): 
    data_path = f'data/{city_name}.csv'
    train_models_path = f'{TRAINED_MODELS_DIR}/{city_name}.joblib'

    # Load the weather data into a pandas DataFrame
    df = pd.read_csv(data_path)

    # Drop columns that contain non-numeric data
    df = df.drop(['city_resolvedAddress', 'city_address', 'city_timezone', 'city_tzoffset', 'day_conditions', 'day_description', 'day_icon', 'day_source', 'day_preciptype', 'day_stations', 'hour_datetime', 'hour_datetimeEpoch', 'hour_preciptype', 'hour_icon', 'hour_source', 'hour_stations', 'day_datetime'], axis=1)

    # Convert hour_conditions to numerical values using label encoding
    label_encoder = LabelEncoder()
    df['hour_conditions'] = label_encoder.fit_transform(df['hour_conditions'])

    # Extract hour, minute, and second values from the day_sunrise column
    df['sunrise_hour'] = df['day_sunrise'].str.split(':').str[0].astype(int)
    df['sunrise_minute'] = df['day_sunrise'].str.split(':').str[1].astype(int)
    df['sunrise_second'] = df['day_sunrise'].str.split(':').str[2].astype(int)

    # Drop the original day_sunrise column
    df = df.drop(['day_sunrise'], axis=1)

    # Extract hour, minute, and second values from the day_sunset column
    df['sunset_hour'] = df['day_sunset'].str.split(':').str[0].astype(int)
    df['sunset_minute'] = df['day_sunset'].str.split(':').str[1].astype(int)
    df['sunset_second'] = df['day_sunset'].str.split(':').str[2].astype(int)

    # Drop the original day_sunset column
    df = df.drop(['day_sunset'], axis=1)

    # Fill missing values with mean values of the column
    imputer = SimpleImputer()
    df = pd.DataFrame(imputer.fit_transform(df), columns=df.columns)

    # Split the data into training and testing sets
    X_train, X_test, y_train, y_test = train_test_split(df.drop('hour_conditions', axis=1), df['hour_conditions'], test_size=0.2, random_state=42)

    # Train a random forest regression model
    try:
        model = joblib.load(train_models_path)
    except:
        model = RandomForestRegressor(n_estimators=100, random_state=42)
        model.fit(X_train, y_train)
        
        # Save the trained model to a file
        joblib.dump(model, train_models_path)


    # Predict weather conditions for the next twelve hours
    new_data = df.tail(12).drop('hour_conditions', axis=1)
    predictions = model.predict(new_data)

    # Inverse transform the label encoding to get the predicted weather conditions
    predictions = label_encoder.inverse_transform(predictions.astype(int))

    # Print the predicted weather conditions
    print(predictions)

def train_all_models():
    # Open the CSV file and create a reader object
    with open(f'{WEATHER_DATA}', 'r') as file:
        reader = csv.DictReader(file)

        # Create a dictionary to store the data for each city
        city_data = {}

        # Iterate through each row of the CSV file
        for row in reader:
            # Get the name of the city for the current row
            city_name = row['city_address']

            # If this is the first time we've seen this city, create a new list for its data
            if city_name not in city_data:
                city_data[city_name] = []

            # Add the current row's data to the list for this city
            city_data[city_name].append(row)

    # Iterate through each city's data and write it to a separate CSV file
    for city_name, data in city_data.items():
        with open(f'{DATA_DIR}/{city_name}.csv', 'w', newline='') as file:
            writer = csv.DictWriter(file, fieldnames=reader.fieldnames)
            writer.writeheader()
            writer.writerows(data)
        train_models(city_name)

if __name__ == "__main__":
    if not os.path.exists(TRAINED_MODELS_DIR):
        os.mkdir(TRAINED_MODELS_DIR)
    if not os.path.exists(DATA_DIR):
        os.mkdir(DATA_DIR)
        train_all_models()