from flask import Flask, request, jsonify
from flask_cors import CORS
import datetime

app = Flask(__name__)
CORS(app)

guard_analysis_history = []

@app.route('/analyze', methods=['POST'])
def analyze_threats():
    data = request.json
    if not data:
        return jsonify({"msg": "missing data"}), 400
    
    threat_level = 0
    raw_message = data.get('message', '').lower()
    
    if any(word in raw_message for word in ['attack', 'exploit', 'sql', 'bypass', 'admin', 'root', 'login', 'password']):
        threat_level += 5
    if any(word in raw_message for word in ['curl', 'wget', 'chmod', 'rm -rf']):
        threat_level += 3
    if data.get('level', 0) > 3:
        threat_level += 2
        
    result = {
        "time": datetime.datetime.now().isoformat(),
        "score": threat_level,
        "note": "Critical" if threat_level > 4 else "Safe",
        "by": "Aether Guard"
    }
    
    guard_analysis_history.append(result)
    return jsonify(result)

@app.route('/history', methods=['GET'])
def get_history():
    return jsonify(guard_analysis_history)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, threaded=True)
