from escpos.printer import Usb, Dummy

from flask import Flask, request
from flask_restful import Resource, Api
from json import dumps
from flask_jsonpify import jsonify
from flask_cors import CORS
import json

from threading import Lock

mutex = Lock()

app = Flask(__name__)
CORS(app)
api = Api(app)

p = Usb(0x04b8, 0x0e15, 0)

class Print(Resource):
    def post(self):
        mutex.acquire()
        try:
            content = request.json
            if content == None:
                return {'status': 'error', 'error': 'No Content'}
            
            p.set(width=4, height=4)
            p.text(str(content['ticketCount'])+"\n")
            p.set(width=2, height=2)
            p.text(content['doelgroepnummer']+"\n")
            p.text(content['naam'] + " " + content['voornaam']+"\n")
            if content['typeVoeding'] is not None: 
                if content['typeVoeding'] != "gewoon": 
                    p.set(align='right',width=2, height=2)
                p.text(content['typeVoeding']+"\n")
                if content['typeVoeding'] != "gewoon": 
                    p.set(align='left',width=2, height=2)
            if content['code'] is not None: 
                p.text(content['code']+"\n")
            if content['volwassenen'] is not None: 
                p.text("Volwassenen: " + str(content['volwassenen'])+"\n")
            if content['kinderen'] is not None: 
                p.text("Kinderen: " + str(content['kinderen'])+"\n")
            if content['specialeVoeding'] is not None: 
                p.text("\n"+content['specialeVoeding']+"\n")
            if content['needsVerjaardag']: 
                p.text("\nVERJAARDAG\n")
            if content['needsMelkpoeder']: 
                p.text("\MELKPOEDER\n")
            p.cut()
            return {'status': 'ok'}
        except:
            return {'status': 'error'}
        finally:
            mutex.release()
            pass


class Eenmaligen(Resource):
    def post(self):
        mutex.acquire()
        try:
            content = request.json
            if content == None:
                return {'status': 'error', 'error': 'No Content'}
            
            p.set(width=4, height=4)
            p.text("VR")
            p.set(width=2, height=2)
            p.text(content['eenmaligenNummer']+"\n")
            p.text(content['naam'] +"\n")
            if content['typeVoeding'] is not None: 
                if content['typeVoeding'] != "gewoon": 
                    p.set(align='right',width=2, height=2)
                p.text(content['typeVoeding']+"\n")
                if content['typeVoeding'] != "gewoon": 
                    p.set(align='left',width=2, height=2)
            if content['grootte'] is not None: 
                p.text(content['grootte']+"\n")
            p.cut()
            return {'status': 'ok'}
        except:
            return {'status': 'error'}
        finally:
            mutex.release()
            pass
        

api.add_resource(Print, '/print')
api.add_resource(Eenmaligen, '/eenmaligen')

if __name__ == '__main__':
     app.run(port='8080')
     