# Python
import os

# Keras
os.environ["KERAS_BACKEND"] = "torch"
import keras
from keras.models import Sequential, Model
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
