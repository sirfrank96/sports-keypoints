from keras.models import load_model

class GolfMLInferrer():
    def __init__(self, cvclient, root_dir):
        self.cvclient = cvclient
        self.root_dir = root_dir

    def predict_is_swing(input):
        is_swing_model = load_model("models/classify_swing_model_openposedata_batch16_epochs50_early_stop.keras")
        is_swing_model.predict(input)
        return