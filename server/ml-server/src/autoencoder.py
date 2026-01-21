# Python
import os
import numpy as np

# Keras
os.environ["KERAS_BACKEND"] = "torch"
from keras.callbacks import EarlyStopping
from keras.models import Sequential, Model, load_model
from keras.layers import LSTM, RepeatVector, TimeDistributed, Dense

def autoencode(x_train, y_train, x_val, y_val):
    encoding_dim = 32 # dimension of compressed projected space
    num_frames = 150 # number of video frames
    num_features = 25*3 # number of datapoints

    model = Sequential()

    # Encoder
    model.add(LSTM(encoding_dim, activation='relu', input_shape=(num_frames, num_features), return_sequences=False, name='encoding_layer'))
    #model.add(LSTM(encoding_dim, activation='relu', return_sequences=False))
    model.add(RepeatVector(num_frames)) # Repeats the encoded vector for the decoder

    # Decoder
    model.add(LSTM(num_features, activation='relu', return_sequences=True, name='decoding_layer'))
    model.add(TimeDistributed(Dense(num_features, activation='sigmoid'))) # Ensures output matches input shape

    model.compile(optimizer='adam', loss='mean_squared_error')

    res = model.fit(x_train, y_train, epochs=10, batch_size=2, shuffle=True)
    print(f'res: {res}, model summary: {model.summary()}')

    intermediate_layer_model = Model(inputs=model.inputs, outputs=model.get_layer('encoding_layer').output)
    intermediate_output = intermediate_layer_model.predict(x_val)
    print(f'intermediate output shape: {intermediate_output.shape}, intermediate output: {intermediate_output}')

def encode_and_classify(x_train, y_train, x_val, y_val):
    x_train_np = np.array(x_train, dtype=float)
    y_train_np = np.array(y_train, dtype=float)
    x_val_np = np.array(x_val, dtype=float)
    y_val_np = np.array(y_val, dtype=float)
    encoding_dim = 32 # dimension of compressed projected space
    num_frames = 150 # number of video frames
    num_features = 25*3 # number of datapoints

    model = Sequential()
    # Encoder
    model.add(LSTM(encoding_dim, activation='relu', input_shape=(num_frames, num_features), return_sequences=True, name='encoding_layer'))
    #model.add(RepeatVector(num_frames)) # Repeats the encoded vector for the decoder

    # Binary classification (e.g., 0 or 1)
    model.add(Dense(1, activation='sigmoid'))

    early_stopping = EarlyStopping(monitor='val_loss', patience=5, restore_best_weights=True)
    loss_function = 'binary_crossentropy'
    metrics_list = ['accuracy']
    model.compile(optimizer='adam', loss=loss_function, metrics=metrics_list)

    model.fit(x_train_np, y_train_np, epochs=50, batch_size=16, shuffle=True, callbacks=[early_stopping])

    res = model.predict(x_val_np)
    print(f"result for predicting x: {x_val_np} is {res}, actual is {y_val_np}")
    print(f"dimensions: x: {x_val_np.shape}, y: {y_val_np.shape}, res: {res.shape}")

    print(f"result for predicting the 3rd video: {x_val_np[3]} is {res[3]}, actual is {y_val_np[3]}")

    print(f"result for predicting the 8th video: {x_val_np[8]} is {res[8]}, actual is {y_val_np[8]}")

    model.save("models/classify_faceon_model_batch16_epochs50_early_stop.keras")

def load_faceon_classification_model():
    model = load_model("models/classify_faceon_model.keras")


