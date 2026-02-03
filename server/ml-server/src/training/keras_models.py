# Python
import os

# Keras
os.environ["KERAS_BACKEND"] = "torch"
from keras.callbacks import EarlyStopping
from keras.models import Sequential, Model
from keras.layers import LSTM, RepeatVector, TimeDistributed, Dense, Masking, Input, GlobalAveragePooling2D
from keras.applications import MobileNetV3Small

# TODO: Add masking layer? and dropouts?
def create_sequential_autoencoder(x_train, y_train, x_val, y_val, num_frames, num_features):
    encoding_dim = 32 # dimension of compressed projected space
    model = Sequential()
    # Encoder
    model.add(LSTM(encoding_dim, activation='relu', input_shape=(num_frames, num_features), return_sequences=False, name='encoding_layer'))
    model.add(RepeatVector(num_frames)) # Repeats the encoded vector for the decoder
    # Decoder
    model.add(LSTM(num_features, activation='relu', return_sequences=True, name='decoding_layer'))
    model.add(TimeDistributed(Dense(num_features, activation='sigmoid'))) # Ensures output matches input shape
    # Compile and train model
    early_stopping = EarlyStopping(monitor='loss', patience=10, restore_best_weights=True)
    metrics_list = ['accuracy']
    model.compile(optimizer='adam', loss='mean_squared_error', metrics=metrics_list)
    res = model.fit(x_train, y_train, epochs=50, batch_size=16, shuffle=True, validation_data=(x_val, y_val), callbacks=[early_stopping])
    print(f'res: {res}, model summary: {model.summary()}')
    return model

def get_encoding_layer_from_autoencoder(autoencoder):
    intermediate_layer_model = Model(inputs=autoencoder.inputs, outputs=autoencoder.get_layer('encoding_layer').output)
    return intermediate_layer_model

def extract_features_using_encoder(encoder, input):
    return encoder.predict(input)

# TODO: Add masking layer? and dropouts?
def create_sequential_classifier(x_train, y_train, x_val, y_val, num_frames, num_features):
    encoding_dim = 32 # dimension of compressed projected space
    model = Sequential()
    # Encoder
    model.add(LSTM(encoding_dim, activation='relu', input_shape=(num_frames, num_features), return_sequences=True, name='encoding_layer'))
    # Binary classification layer (e.g., 0 or 1)
    model.add(Dense(1, activation='sigmoid'))
    # Compile and train model
    early_stopping = EarlyStopping(monitor='loss', patience=5, restore_best_weights=True)
    loss_function = 'binary_crossentropy'
    metrics_list = ['accuracy']
    model.compile(optimizer='adam', loss=loss_function, metrics=metrics_list)
    res = model.fit(x_train, y_train, epochs=50, batch_size=16, shuffle=True, validation_data=(x_val, y_val), callbacks=[early_stopping])
    print(f'res: {res}, model summary: {model.summary()}')
    return model

def predict_using_sequential_classifier(classifier, input):
    return classifier.predict(input)

# TODO: Add masking layer? and dropouts?
def create_sequential_cnn_classifier(x_train, y_train, x_val, y_val, num_frames, num_features):
    encoding_dim = 32 # dimension of compressed projected space
    model = Sequential()
    # CNN extract features from frame images (uses google's MobileNetV2)
    mobile_net_model = MobileNetV3Small(include_top=False, weights='imagenet', pooling='avg', input_shape=(256, 160, 3), include_preprocessing=True)
    mobile_net_model.trainable = False
    model.add(TimeDistributed(mobile_net_model))
    # LSTM
    sequence_model = LSTM(encoding_dim, activation='relu', return_sequences=True, name='encoding_layer')
    model.add(sequence_model)
    # Binary classification layer (e.g., 0 or 1)
    model.add(Dense(1, activation='sigmoid'))
    # Compile and train model
    early_stopping = EarlyStopping(monitor='loss', patience=10, restore_best_weights=True)
    loss_function = 'binary_crossentropy'
    metrics_list = ['accuracy']
    model.compile(optimizer='adam', loss=loss_function, metrics=metrics_list)
    res= model.fit(x_train, y_train, epochs=50, batch_size=16, shuffle=True, validation_data=(x_val, y_val), callbacks=[early_stopping])
    print(f'res: {res}, model summary: {model.summary()}')
    return model

def predict_using_sequential_cnn_classifier(seq_classifier, input):
    return seq_classifier.predict(input)

#x_train_np = np.array(x_train, dtype=float)
#y_train_np = np.array(y_train, dtype=float)
#x_val_np = np.array(x_val, dtype=float)
#y_val_np = np.array(y_val, dtype=float) 

#print(f"result for predicting input: {input} is {res}, actual is {y_val}")
#print(f"dimensions: x: {x_val.shape}, y: {y_val.shape}, res: {res.shape}")
#print(f"result for predicting the 3rd video: {x_val[3]} is {res[3]}, actual is {y_val[3]}")
#print(f"result for predicting the 8th video: {x_val[7]} is {res[7]}, actual is {y_val[7]}")   

# encode_and_classify_swing_motion
# identifies in video what is motion and what is just setup or end of video
# split into chunks? (setup, takeaway, backswing, downswing, impact, followthrough)
# when there is consecutive movement, this is the start, when there is consecutive non movement, this is not



